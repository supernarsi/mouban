package dao

import (
	"mouban/common"
	"mouban/consts"
	"mouban/model"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var (
	dataProcessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mouban_data_process_total",
		Help: "Data processed counter",
	}, []string{"type", "result"})
)

func GetSchedule(doubanId uint64, t uint8) *model.Schedule {
	schedule := &model.Schedule{}
	result := common.Db.Where("douban_id = ? AND type = ?", doubanId, t).Find(schedule)
	if result.Error != nil {
		logrus.Errorln("get schedule error:", result.Error, "douban_id:", doubanId, "type:", t)
		return nil
	}
	if schedule.ID == 0 {
		return nil
	}
	return schedule
}

// SearchScheduleByStatus idx_status
func SearchScheduleByStatus(t uint8, status uint8) *model.Schedule {
	schedule := &model.Schedule{}
	result := common.Db.Where("type = ? AND status = ?", t, status).
		Order("updated_at asc").
		Limit(1).
		Find(&schedule)
	if result.Error != nil {
		logrus.Errorln("search schedule by status error:", result.Error, "type:", t, "status:", status)
		return nil
	}
	if schedule.ID == 0 {
		return nil
	}
	return schedule
}

// SearchScheduleByAll idx_search
func SearchScheduleByAll(t uint8, status uint8, result uint8) *model.Schedule {
	schedule := &model.Schedule{}
	dbResult := common.Db.Where("type = ? AND `status`= ? AND result = ?", t, status, result).
		Order("updated_at asc").
		Limit(1).
		Find(&schedule)
	if dbResult.Error != nil {
		logrus.Errorln("search schedule by all error:", dbResult.Error, "type:", t, "status:", status, "result:", result)
		return nil
	}
	if schedule.ID == 0 {
		return nil
	}
	return schedule
}

// CasOrphanSchedule idx_status
func CasOrphanSchedule(t uint8, expire time.Duration) int64 {
	result := common.Db.Model(&model.Schedule{}).
		Where("type = ? AND status = ? AND updated_at < ?", t, consts.ScheduleCrawling.Code, time.Now().Add(-expire)).
		Update("status", consts.ScheduleToCrawl.Code)
	if result.Error != nil {
		logrus.Errorln("cas orphan schedule error:", result.Error, "type:", t)
		return 0
	}
	return result.RowsAffected
}

// CasScheduleStatus uk_schedule
func CasScheduleStatus(doubanId uint64, t uint8, status uint8, rawStatus uint8) bool {
	result := common.Db.Model(&model.Schedule{}).
		Where("douban_id = ? AND type = ? AND status = ?", doubanId, t, rawStatus).
		Update("status", status)
	if result.Error != nil {
		logrus.Errorln("cas schedule status error:", result.Error, "douban_id:", doubanId, "type:", t)
		return false
	}
	return result.RowsAffected > 0
}

// ChangeScheduleResult uk_schedule
func ChangeScheduleResult(doubanId uint64, t uint8, result uint8) {
	dataProcessTotal.WithLabelValues(consts.ParseType(t).Name, consts.ParseResult(result).Name).Inc()

	dbResult := common.Db.Model(&model.Schedule{}).
		Where("douban_id = ? AND type = ?", doubanId, t).
		Update("result", result)
	if dbResult.Error != nil {
		logrus.Errorln("change schedule result error:", dbResult.Error, "douban_id:", doubanId, "type:", t)
	}
}

func CreateScheduleNx(doubanId uint64, t uint8, status uint8, result uint8) bool {
	data := &model.Schedule{}
	insert := &model.Schedule{
		DoubanId: doubanId,
		Type:     t,
		Status:   &status,
		Result:   &result,
	}
	dbResult := common.Db.Where("douban_id = ? AND type = ?", doubanId, t).Attrs(insert).FirstOrCreate(data)
	if dbResult.Error != nil {
		logrus.Errorln("create schedule nx error:", dbResult.Error, "douban_id:", doubanId, "type:", t)
		return false
	}
	return dbResult.RowsAffected > 0
}
