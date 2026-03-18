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
				logrus.WithField("douban_id", pendingSchedule.DoubanId).Debugln("pending", t.Name, "item found")

				changed := dao.CasScheduleStatus(pendingSchedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *pendingSchedule.Status)
				if changed {
					ch <- pendingSchedule
					continue
				}
			}

			retrySchedule := dao.SearchScheduleByAll(t.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
			if retrySchedule != nil {
				logrus.WithField("douban_id", retrySchedule.DoubanId).Debugln("retry", t.Name, "item found")
				changed := dao.CasScheduleStatus(retrySchedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *retrySchedule.Status)
				if changed {
					ch <- retrySchedule
					continue
				}
			}

			discoverSchedule := dao.SearchScheduleByStatus(t.Code, consts.ScheduleCanCrawl.Code)
			if discoverSchedule != nil {
				logrus.WithField("douban_id", discoverSchedule.DoubanId).Debugln("discover", t.Name, "item found")
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
		logrus.WithFields(logrus.Fields{
			"type":       t.Name,
			"douban_id":  strconv.FormatUint(schedule.DoubanId, 10),
			"thread_idx": index,
		}).Infoln("item processing started")
		processItem(schedule.Type, schedule.DoubanId)
		dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
		logrus.WithFields(logrus.Fields{
			"type":       t.Name,
			"douban_id":  strconv.FormatUint(schedule.DoubanId, 10),
			"thread_idx": index,
		}).Infoln("item processing completed")
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

	logrus.WithField("concurrency", concurrency).Infoln("item agent(s) enabled")
}
