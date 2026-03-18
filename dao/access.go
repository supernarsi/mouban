package dao

import (
	"mouban/common"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func AddAccess(doubanUid uint64, path string, ip string, ua string, referer string) {
	access := &model.Access{
		DoubanUid: doubanUid,
		Path:      path,
		Ip:        ip,
		UserAgent: ua,
		Referer:   referer,
	}

	result := common.Db.Create(access)
	if result.Error != nil {
		logrus.Errorln("add access error:", result.Error, "douban_uid:", doubanUid, "path:", path)
	}
}
