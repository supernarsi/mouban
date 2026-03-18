package dao

import (
	"mouban/common"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func CountGame() int64 {
	var count int64
	result := common.Db.Model(&model.Game{}).Count(&count)
	if result.Error != nil {
		logrus.Errorln("count game error:", result.Error)
		return 0
	}
	return count
}

func UpsertGame(game *model.Game) {
	logrus.WithField("upsert", "game").Infoln("upsert game", game.DoubanId, game.Title)
	data := &model.Game{}
	result := common.Db.Where("douban_id = ?", game.DoubanId).Assign(game).FirstOrCreate(data)
	if result.Error != nil {
		logrus.Errorln("upsert game error:", result.Error, "douban_id:", game.DoubanId)
	}
}

func UpdateGameThumbnail(doubanId uint64, thumbnail string) {
	result := common.Db.Model(&model.Game{}).Where("douban_id = ?", doubanId).Update("thumbnail", thumbnail)
	if result.Error != nil {
		logrus.Errorln("update game thumbnail error:", result.Error, "douban_id:", doubanId)
	}
}

func CreateGameNx(game *model.Game) bool {
	data := &model.Game{}
	inserted := common.Db.Where("douban_id = ?", game.DoubanId).Attrs(game).FirstOrCreate(data).RowsAffected > 0
	if inserted {
		logrus.Infoln("create game", game.DoubanId, game.Title)
	}
	return inserted
}

func GetGameDetail(doubanId uint64) *model.Game {
	game := &model.Game{}
	result := common.Db.Where("douban_id = ?", doubanId).Find(game)
	if result.Error != nil {
		logrus.Errorln("get game detail error:", result.Error, "douban_id:", doubanId)
		return nil
	}
	if game.ID == 0 {
		return nil
	}
	return game
}

func ListGameBrief(doubanIds *[]uint64) *[]model.Game {
	var games *[]model.Game
	result := common.Db.Omit("intro").Where("douban_id IN ?", *doubanIds).Find(&games)
	if result.Error != nil {
		logrus.Errorln("list game brief error:", result.Error)
		return nil
	}
	return games
}
