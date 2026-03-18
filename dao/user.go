package dao

import (
	"mouban/common"
	"mouban/model"
	"time"

	"github.com/sirupsen/logrus"
)

func CountUser() int64 {
	var count int64
	result := common.Db.Model(&model.User{}).Count(&count)
	if result.Error != nil {
		logrus.Errorln("count user error:", result.Error)
		return 0
	}
	return count
}

func UpsertUser(user *model.User) {
	logrus.WithField("upsert", "user").Infoln("upsert user", user.DoubanUid, user.Name)
	data := &model.User{}
	result := common.Db.Where("douban_uid = ?", user.DoubanUid).Assign(user).FirstOrCreate(data)
	if result.Error != nil {
		logrus.Errorln("upsert user error:", result.Error, "douban_uid:", user.DoubanUid)
	}
}

func RefreshUser(user *model.User) {
	logrus.Infoln("refresh user", user.DoubanUid, user.Name)
	result := common.Db.Model(&model.User{}).
		Where("douban_uid = ?", user.DoubanUid).
		Updates(model.User{CheckAt: time.Unix(0, 0), SyncAt: time.Unix(0, 0), PublishAt: time.Unix(0, 0)})
	if result.Error != nil {
		logrus.Errorln("refresh user error:", result.Error, "douban_uid:", user.DoubanUid)
	}
}

func GetUser(doubanUid uint64) *model.User {
	if doubanUid == 0 {
		return nil
	}
	user := &model.User{}
	result := common.Db.Where("douban_uid = ?", doubanUid).Find(user)
	if result.Error != nil {
		logrus.Errorln("get user error:", result.Error, "douban_uid:", doubanUid)
		return nil
	}
	if user.ID == 0 {
		return nil
	}
	return user
}

func GetUserByDomain(domain string) *model.User {
	if domain == "" {
		return nil
	}
	user := &model.User{}
	result := common.Db.Where("domain = ?", domain).Find(user)
	if result.Error != nil {
		logrus.Errorln("get user by domain error:", result.Error, "domain:", domain)
		return nil
	}
	if user.ID == 0 {
		return nil
	}
	return user
}
