package agent

import (
	"mouban/consts"
	"mouban/crawl"
	"mouban/dao"
	"mouban/model"
	"mouban/util"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// syncComment 通用函数：同步用户评论和条目数据
// crawlFn: 爬取评论和条目的函数
// createItemFn: 创建条目的 DAO 函数（接收指针）
// createScheduleFn: 创建调度任务的 DAO 函数
// typeName: 类型名称（用于日志）
// getDoubanId: 获取条目 ID 的函数
func syncComment[T any](
	user *model.User,
	forceSyncAfter time.Time,
	crawlFn func(*model.User, time.Time) (*[]model.Comment, *[]T, error),
	createItemFn func(*T) bool,
	createScheduleFn func(uint64, uint8, uint8, uint8) bool,
	typeName string,
	typeCode uint8,
	getDoubanId func(*T) uint64,
) {
	comment, items, err := crawlFn(user, forceSyncAfter)
	if err != nil {
		panic(err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"type": typeName,
				}).Errorln("sync comment panic", r, "=>", util.GetCurrentGoroutineStack())
			}
		}()

		// 如果是全量同步，隐藏已删除的评论
		if forceSyncAfter.Unix() == 0 {
			newCommentIds := make(map[uint64]bool)
			for i := range *items {
				newCommentIds[getDoubanId(&(*items)[i])] = true
			}
			oldCommentIds := dao.GetCommentIds(user.DoubanUid, typeCode)
			for i := range *oldCommentIds {
				id := (*oldCommentIds)[i]
				if !newCommentIds[id] {
					dao.HideComment(user.DoubanUid, typeCode, id)
				}
			}
		}

		// 处理每条评论和条目
		for i := range *items {
			dao.UpsertComment(&(*comment)[i])

			item := &(*items)[i]
			added := createItemFn(item)
			if added {
				createScheduleFn(getDoubanId(item), typeCode, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
			}
		}
	}()
}

func processItem(t uint8, doubanId uint64) {
	switch t {
	case consts.TypeBook.Code:
		processBook(doubanId)
	case consts.TypeMovie.Code:
		processMovie(doubanId)
	case consts.TypeGame.Code:
		processGame(doubanId)
	case consts.TypeSong.Code:
		processSong(doubanId)
	}
}

func processBook(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("process book panic", doubanId, r, "=>", util.GetCurrentGoroutineStack())
		}
	}()
	book, rating, newUser, newItems, err := crawl.Book(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeBook)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeBook.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}

	dao.UpsertBook(book)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeBook.Code, consts.ScheduleReady.Code)
}

func processMovie(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("process movie panic", doubanId, r, "=>", util.GetCurrentGoroutineStack())
		}
	}()
	movie, rating, newUser, newItems, err := crawl.Movie(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeMovie)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeMovie.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}
	dao.UpsertMovie(movie)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeMovie.Code, consts.ScheduleReady.Code)
}

func processGame(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("process game panic", doubanId, r, "=>", util.GetCurrentGoroutineStack())
		}
	}()

	game, rating, newUser, newItems, err := crawl.Game(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeGame)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeGame.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}
	dao.UpsertGame(game)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeGame.Code, consts.ScheduleReady.Code)
}

func processSong(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("process song panic", doubanId, r, "=>", util.GetCurrentGoroutineStack())
		}
	}()

	song, rating, newUser, newItems, err := crawl.Song(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeSong)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeSong.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}
	dao.UpsertSong(song)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeSong.Code, consts.ScheduleReady.Code)
}

func processDiscoverUser(newUsers *[]string) {
	if newUsers == nil {
		return
	}
	level := viper.GetInt("agent.discover.level")
	if level == 0 {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln("process discover user panic", r, "=>", util.GetCurrentGoroutineStack())
			}
		}()
		totalFound := len(*newUsers)
		newFound := 0
		for _, idOrDomain := range *newUsers {
			id, err := strconv.ParseUint(idOrDomain, 10, 64)
			if err != nil {
				if level == 1 {
					continue
				}
				user := dao.GetUserByDomain(idOrDomain)
				if user == nil {
					id = crawl.UserId(idOrDomain)
				}
			}
			if id > 0 {
				added := dao.CreateScheduleNx(id, consts.TypeUser.Code, consts.ScheduleCanCrawl.Code, consts.ScheduleUnready.Code)
				if added {
					newFound += 1
				}
			}
		}
		if newFound > 0 {
			logrus.WithFields(logrus.Fields{
				"new_found":  newFound,
				"total_found": totalFound,
			}).Infoln("users discovered")
		}
	}()
}

func processDiscoverItem(newItems *[]uint64, t consts.Type) {
	if newItems == nil || len(*newItems) == 0 {
		return
	}
	level := viper.GetInt("agent.discover.level")
	if level == 0 {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln("process discover item panic", r, "=>", util.GetCurrentGoroutineStack())
			}
		}()
		totalFound := len(*newItems)
		newFound := 0
		for _, doubanId := range *newItems {
			added := dao.CreateScheduleNx(doubanId, t.Code, consts.ScheduleCanCrawl.Code, consts.ScheduleUnready.Code)
			if added {
				newFound += 1
			}
		}
		if newFound > 0 {
			logrus.WithFields(logrus.Fields{
				"new_found":  newFound,
				"total_found": totalFound,
				"type":       t.Name,
			}).Infoln("items discovered")
		}
	}()
}

func processUser(doubanUid uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("process user panic", doubanUid, r, "=>", util.GetCurrentGoroutineStack())
		}
	}()

	userPublish, _ := crawl.UserPublish(doubanUid)
	rawUser := dao.GetUser(doubanUid)
	if rawUser != nil && rawUser.PublishAt.Equal(userPublish) {
		logrus.WithField("douban_uid", doubanUid).Debugln("user not changed")
		rawUser.CheckAt = time.Now()
		dao.UpsertUser(rawUser)
		return
	}
	logrus.WithField("douban_uid", doubanUid).Infoln("user changed")

	//user
	user, err := crawl.UserOverview(doubanUid)
	if err != nil {
		dao.ChangeScheduleResult(doubanUid, consts.TypeUser.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}

	// choose update type
	forceSyncAfter := time.Unix(0, 0)
	if rawUser != nil && rawUser.SyncAt.AddDate(1, 0, 0).After(time.Now()) {
		forceSyncAfter = rawUser.SyncAt
	}

	logrus.WithFields(logrus.Fields{
		"douban_uid":   doubanUid,
		"sync_after":   forceSyncAfter,
	}).Infoln("user sync timestamp")

	//book
	if user.BookDo+user.BookWish+user.BookCollect > 0 {
		syncCommentBook(user, forceSyncAfter)
	}

	//movie
	if user.MovieDo+user.MovieWish+user.MovieCollect > 0 {
		syncCommentMovie(user, forceSyncAfter)
	}

	//game
	if user.GameDo+user.GameWish+user.GameCollect > 0 {
		syncCommentGame(user, forceSyncAfter)
	}

	//song
	if user.SongDo+user.SongWish+user.SongCollect > 0 {
		syncCommentSong(user, forceSyncAfter)
	}

	user.CheckAt = time.Now()
	user.SyncAt = time.Now()

	dao.UpsertUser(user)
	dao.ChangeScheduleResult(doubanUid, consts.TypeUser.Code, consts.ScheduleReady.Code)
}

func syncCommentGame(user *model.User, forceSyncAfter time.Time) {
	syncComment(
		user,
		forceSyncAfter,
		crawl.CommentGame,
		dao.CreateGameNx,
		dao.CreateScheduleNx,
		"game",
		consts.TypeGame.Code,
		func(item *model.Game) uint64 { return item.DoubanId },
	)
}

func syncCommentBook(user *model.User, forceSyncAfter time.Time) {
	syncComment(
		user,
		forceSyncAfter,
		crawl.CommentBook,
		dao.CreateBookNx,
		dao.CreateScheduleNx,
		"book",
		consts.TypeBook.Code,
		func(item *model.Book) uint64 { return item.DoubanId },
	)
}

func syncCommentMovie(user *model.User, forceSyncAfter time.Time) {
	syncComment(
		user,
		forceSyncAfter,
		crawl.CommentMovie,
		dao.CreateMovieNx,
		dao.CreateScheduleNx,
		"movie",
		consts.TypeMovie.Code,
		func(item *model.Movie) uint64 { return item.DoubanId },
	)
}

func syncCommentSong(user *model.User, forceSyncAfter time.Time) {
	syncComment(
		user,
		forceSyncAfter,
		crawl.CommentSong,
		dao.CreateSongNx,
		dao.CreateScheduleNx,
		"song",
		consts.TypeSong.Code,
		func(item *model.Song) uint64 { return item.DoubanId },
	)
}
