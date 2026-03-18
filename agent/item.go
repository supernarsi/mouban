package agent

import (
	"mouban/consts"
	"mouban/dao"
	"mouban/model"
	"mouban/util"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func itemSelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("item selector panic", r, util.GetCurrentGoroutineStack())
		}
	}()

	types := []consts.Type{consts.TypeBook, consts.TypeMovie, consts.TypeGame, consts.TypeSong}
	for {
		for _, t := range types {
			pendingSchedule := dao.SearchScheduleByStatus(t.Code, consts.ScheduleToCrawl.Code)
			if pendingSchedule != nil {
				logrus.Infoln("pending", t.Name, "item found", pendingSchedule.DoubanId)

				changed := dao.CasScheduleStatus(pendingSchedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *pendingSchedule.Status)
				if changed {
					ch <- pendingSchedule
					continue
				}
			}

			retrySchedule := dao.SearchScheduleByAll(t.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
			if retrySchedule != nil {
				logrus.Infoln("retry", t.Name, "item found", retrySchedule.DoubanId)
				changed := dao.CasScheduleStatus(retrySchedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *retrySchedule.Status)
				if changed {
					ch <- retrySchedule
					continue
				}
			}

			discoverSchedule := dao.SearchScheduleByStatus(t.Code, consts.ScheduleCanCrawl.Code)
			if discoverSchedule != nil {
				logrus.Infoln("discover", t.Name, "item found", discoverSchedule.DoubanId)
				changed := dao.CasScheduleStatus(discoverSchedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *discoverSchedule.Status)
				if changed {
					ch <- discoverSchedule
					continue
				}
			}
		}
		time.Sleep(time.Minute)
	}
}

func itemWorker(index int, ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("item worker panic", r, "item worker (", index, ") crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	for schedule := range ch {
		t := consts.ParseType(schedule.Type)
		logrus.Infoln("item thread", index, "start", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
		processItem(schedule.Type, schedule.DoubanId)
		dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
		logrus.Infoln("item thread", index, "end", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("item agent disabled")
		return
	}

	concurrency := viper.GetInt("agent.item.concurrency")

	ch := make(chan *model.Schedule, concurrency)

	go func() {
		itemSelector(ch)
	}()

	for i := 0; i < concurrency; i++ {
		j := i + 1
		go func() {
			itemWorker(j, ch)
		}()
	}

	logrus.Infoln(concurrency, "item agent(s) enabled")
}
