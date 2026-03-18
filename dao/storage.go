package dao

import (
	"mouban/common"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func UpsertStorage(storage *model.Storage) {
	logrus.WithField("upsert", "storage").Infoln("upsert storage", storage.Source, storage.Target)
	data := &model.Storage{}
	result := common.Db.Where("source = ?", storage.Source).Assign(storage).FirstOrCreate(data)
	if result.Error != nil {
		logrus.Errorln("upsert storage error:", result.Error, "source:", storage.Source)
	}
}

func GetStorageByMd5(md5 string) *model.Storage {
	if md5 == "" {
		return nil
	}
	storage := &model.Storage{}
	result := common.Db.Where("md5 = ?", md5).Limit(1).Find(storage)
	if result.Error != nil {
		logrus.Errorln("get storage by md5 error:", result.Error, "md5:", md5)
		return nil
	}
	if storage.ID == 0 {
		return nil
	}
	return storage
}

func GetStorage(source string) *model.Storage {
	if source == "" {
		return nil
	}
	storage := &model.Storage{}
	result := common.Db.Where("source = ?", source).Find(storage)
	if result.Error != nil {
		logrus.Errorln("get storage error:", result.Error, "source:", source)
		return nil
	}
	if storage.ID == 0 {
		return nil
	}
	return storage
}
