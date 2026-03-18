package dao

import (
	"mouban/common"
	"mouban/consts"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func HideComment(doubanUid uint64, t uint8, doubanId uint64) {
	logrus.Infoln("hide comment for", doubanUid, "type", t, "at", doubanId)
	result := common.Db.Model(&model.Comment{}).
		Where("douban_uid = ? AND type = ? AND douban_id = ?", doubanUid, t, doubanId).
		Update("action", consts.ActionHide.Code)
	if result.Error != nil {
		logrus.Errorln("hide comment error:", result.Error, "douban_uid:", doubanUid, "type:", t, "douban_id:", doubanId)
	}
}

func GetCommentIds(doubanUid uint64, t uint8) *[]uint64 {
	var doubanIds []uint64
	result := common.Db.Model(&model.Comment{}).Where("douban_uid = ? AND type = ?", doubanUid, t).Select("douban_id").Find(&doubanIds)
	if result.Error != nil {
		logrus.Errorln("get comment ids error:", result.Error, "douban_uid:", doubanUid, "type:", t)
		return nil
	}
	return &doubanIds
}

func UpsertComment(comment *model.Comment) {
	logrus.WithField("upsert", "comment").Infoln("upsert comment", comment.DoubanId, comment.Type, "for", comment.DoubanUid)
	data := &model.Comment{}
	result := common.Db.Where("douban_id = ? AND douban_uid = ? AND type = ?", comment.DoubanId, comment.DoubanUid, comment.Type).Assign(comment).FirstOrCreate(data)
	if result.Error != nil {
		logrus.Errorln("upsert comment error:", result.Error, "douban_id:", comment.DoubanId, "douban_uid:", comment.DoubanUid)
	}
}

func GetComment(doubanId uint64, doubanUid uint64, t uint8) *model.Comment {
	comment := &model.Comment{}
	result := common.Db.Where("douban_id = ? AND douban_uid = ? AND type = ?", doubanId, doubanUid, t).Find(comment)
	if result.Error != nil {
		logrus.Errorln("get comment error:", result.Error, "douban_id:", doubanId, "douban_uid:", doubanUid, "type:", t)
		return nil
	}
	if comment.ID == 0 {
		return nil
	}
	return comment
}

// SearchComment idx_search
func SearchComment(doubanUid uint64, t uint8, action uint8, offset int) *[]model.Comment {
	var comment *[]model.Comment
	result := common.Db.Where("douban_uid = ? AND type = ? AND action = ?", doubanUid, t, action).
		Order("mark_date desc").
		Offset(offset).
		Find(&comment)
	if result.Error != nil {
		logrus.Errorln("search comment error:", result.Error, "douban_uid:", doubanUid, "type:", t, "action:", action)
		return nil
	}
	return comment
}
