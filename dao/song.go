package dao

import (
	"mouban/common"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func CountSong() int64 {
	var count int64
	result := common.Db.Model(&model.Song{}).Count(&count)
	if result.Error != nil {
		logrus.Errorln("count song error:", result.Error)
		return 0
	}
	return count
}

func UpsertSong(song *model.Song) {
	logrus.WithField("upsert", "song").Infoln("upsert song", song.DoubanId, song.Title)
	data := &model.Song{}
	result := common.Db.Where("douban_id = ?", song.DoubanId).Assign(song).FirstOrCreate(data)
	if result.Error != nil {
		logrus.Errorln("upsert song error:", result.Error, "douban_id:", song.DoubanId)
	}
}

func UpdateSongThumbnail(doubanId uint64, thumbnail string) {
	result := common.Db.Model(&model.Song{}).Where("douban_id = ?", doubanId).Update("thumbnail", thumbnail)
	if result.Error != nil {
		logrus.Errorln("update song thumbnail error:", result.Error, "douban_id:", doubanId)
	}
}

func CreateSongNx(song *model.Song) bool {
	data := &model.Song{}
	inserted := common.Db.Where("douban_id = ?", song.DoubanId).Attrs(song).FirstOrCreate(data).RowsAffected > 0
	if inserted {
		logrus.Infoln("create song", song.DoubanId, song.Title)
	}
	return inserted
}

func GetSongDetail(doubanId uint64) *model.Song {
	song := &model.Song{}
	result := common.Db.Where("douban_id = ?", doubanId).Find(song)
	if result.Error != nil {
		logrus.Errorln("get song detail error:", result.Error, "douban_id:", doubanId)
		return nil
	}
	if song.ID == 0 {
		return nil
	}
	return song
}

func ListSongBrief(doubanIds *[]uint64) *[]model.Song {
	var songs *[]model.Song
	result := common.Db.Omit("intro", "track_list").Where("douban_id IN ?", *doubanIds).Find(&songs)
	if result.Error != nil {
		logrus.Errorln("list song brief error:", result.Error)
		return nil
	}
	return songs
}
