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

func userSelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user selector panic", r, util.GetCurrentGoroutineStack())
		}
	}()

	for {
		pendingSchedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleToCrawl.Code)
		if pendingSchedule != nil {
			logrus.WithField("douban_id", pendingSchedule.DoubanId).Debugln("pending user found")
			changed := dao.CasScheduleStatus(pendingSchedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *pendingSchedule.Status)
			if changed {
				ch <- pendingSchedule
				continue
			}
		}

		retrySchedule := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)

		if retrySchedule != nil {
			logrus.WithField("douban_id", retrySchedule.DoubanId).Debugln("retry user found")
			changed := dao.CasScheduleStatus(retrySchedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *retrySchedule.Status)
			if changed {
				ch <- retrySchedule
				continue
			}
		}

		if viper.GetBool("agent.flow.discover") {
			discoverSchedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleCanCrawl.Code)
			if discoverSchedule != nil {
				logrus.WithField("douban_id", discoverSchedule.DoubanId).Debugln("discover user found")
				changed := dao.CasScheduleStatus(discoverSchedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *discoverSchedule.Status)
				if changed {
					ch <- discoverSchedule
					continue
				}
			}
		}

		time.Sleep(time.Minute)
	}
}

func userWorker(index int, ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user worker panic", r, "user worker (", index, ") crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	for schedule := range ch {
		t := consts.ParseType(schedule.Type)
		logrus.WithFields(logrus.Fields{
			"douban_id":  strconv.FormatUint(schedule.DoubanId, 10),
			"thread_idx": index,
		}).Infoln("user processing started")
		processUser(schedule.DoubanId)
		dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
		logrus.WithFields(logrus.Fields{
			"douban_id":  strconv.FormatUint(schedule.DoubanId, 10),
			"thread_idx": index,
		}).Infoln("user processing completed")
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("user agent disabled")
		return
	}

	concurrency := viper.GetInt("agent.user.concurrency")

	ch := make(chan *model.Schedule, concurrency)

	go func() {
		userSelector(ch)
	}()

	for i := 0; i < concurrency; i++ {
		j := i + 1
		go func() {
			userWorker(j, ch)
		}()
	}

	logrus.WithField("concurrency", concurrency).Infoln("user agent(s) enabled")

}
