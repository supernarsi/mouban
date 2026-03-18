package dao

import (
	"mouban/common"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func CountMovie() int64 {
	var count int64
	result := common.Db.Model(&model.Movie{}).Count(&count)
	if result.Error != nil {
		logrus.Errorln("count movie error:", result.Error)
		return 0
	}
	return count
}

func UpsertMovie(movie *model.Movie) {
	logrus.WithField("upsert", "movie").Infoln("upsert movie", movie.DoubanId, movie.Title)
	data := &model.Movie{}
	result := common.Db.Where("douban_id = ?", movie.DoubanId).Assign(movie).FirstOrCreate(data)
	if result.Error != nil {
		logrus.Errorln("upsert movie error:", result.Error, "douban_id:", movie.DoubanId)
	}
}

func UpdateMovieThumbnail(doubanId uint64, thumbnail string) {
	result := common.Db.Model(&model.Movie{}).Where("douban_id = ?", doubanId).Update("thumbnail", thumbnail)
	if result.Error != nil {
		logrus.Errorln("update movie thumbnail error:", result.Error, "douban_id:", doubanId)
	}
}

func CreateMovieNx(movie *model.Movie) bool {
	data := &model.Movie{}
	inserted := common.Db.Where("douban_id = ?", movie.DoubanId).Attrs(movie).FirstOrCreate(data).RowsAffected > 0
	if inserted {
		logrus.Infoln("create movie", movie.DoubanId, movie.Title)
	}
	return inserted
}

func GetMovieDetail(doubanId uint64) *model.Movie {
	movie := &model.Movie{}
	result := common.Db.Where("douban_id = ?", doubanId).Find(movie)
	if result.Error != nil {
		logrus.Errorln("get movie detail error:", result.Error, "douban_id:", doubanId)
		return nil
	}
	if movie.ID == 0 {
		return nil
	}
	return movie
}

func ListMovieBrief(doubanIds *[]uint64) *[]model.Movie {
	var movies *[]model.Movie
	result := common.Db.Omit("intro").Where("douban_id IN ?", *doubanIds).Find(&movies)
	if result.Error != nil {
		logrus.Errorln("list movie brief error:", result.Error)
		return nil
	}
	return movies
}
