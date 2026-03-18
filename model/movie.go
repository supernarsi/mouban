package model

import (
	"time"
)

type Movie struct {
	ID          uint64 `gorm:"comment:自增主键 ID"`
	DoubanId    uint64 `gorm:"not null;uniqueIndex;comment:豆瓣电影 ID"`
	Title       string `gorm:"not null;type:varchar(512);comment:电影名称"`
	Director    string `gorm:"type:varchar(512);comment:导演"`
	Writer      string `gorm:"type:varchar(512);comment:编剧"`
	Actor       string `gorm:"type:varchar(2048);comment:主演 (多人)"`
	Style       string `gorm:"type:varchar(512);comment:类型/风格"`
	Site        string `gorm:"type:varchar(512);comment:官方网站"`
	Country     string `gorm:"type:varchar(512);comment:制片国家/地区"`
	Language    string `gorm:"type:varchar(512);comment:语言"`
	PublishDate string `gorm:"type:varchar(512);comment:上映日期"`
	Episode     uint32 `gorm:"comment:集数 (电视剧)"`
	Duration    uint32 `gorm:"comment:片长 (分钟)"`
	Alias       string `gorm:"type:varchar(512);comment:又名"`
	IMDb        string `gorm:"type:varchar(512);column:imdb;comment:IMDb 链接"`
	Intro       string `gorm:"type:mediumtext;comment:简介"`
	Thumbnail   string `gorm:"type:varchar(512);comment:海报 URL"`
	CreatedAt   time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt   time.Time `gorm:"comment:记录更新时间"`
}

func (Movie) TableName() string {
	return "movie"
}

func (movie Movie) Show() *MovieVO {
	return &MovieVO{
		DoubanId:    movie.DoubanId,
		Title:       movie.Title,
		Style:       movie.Style,
		Director:    movie.Director,
		Writer:      movie.Writer,
		Actor:       movie.Actor,
		PublishDate: movie.PublishDate,
		Alias:       movie.Alias,
		Thumbnail:   movie.Thumbnail,
	}
}

type MovieVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Style       string `json:"style"`
	Director    string `json:"director"`
	Writer      string `json:"writer"`
	Actor       string `json:"actor"`
	PublishDate string `json:"publish_date"`
	Alias       string `json:"alias"`
	Thumbnail   string `json:"thumbnail"`
}
