package model

import (
	"time"
)

type Game struct {
	ID          uint64 `gorm:"comment:自增主键 ID"`
	DoubanId    uint64 `gorm:"not null;uniqueIndex;comment:豆瓣游戏 ID"`
	Title       string `gorm:"not null;type:varchar(512);comment:游戏名称"`
	Platform    string `gorm:"type:varchar(512);comment:平台 (PC/PS5/Switch 等)"`
	Genre       string `gorm:"type:varchar(512);comment:类型"`
	Alias       string `gorm:"type:varchar(512);comment:又名"`
	Developer   string `gorm:"type:varchar(512);comment:开发商"`
	Publisher   string `gorm:"type:varchar(512);comment:发行商"`
	PublishDate string `gorm:"type:varchar(512);comment:发行日期"`
	Intro       string `gorm:"type:mediumtext;comment:游戏简介"`
	Thumbnail   string `gorm:"type:varchar(512);comment:封面图 URL"`
	CreatedAt   time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt   time.Time `gorm:"comment:记录更新时间"`
}

func (Game) TableName() string {
	return "game"
}

func (game Game) Show() *GameVO {
	return &GameVO{
		DoubanId:    game.DoubanId,
		Title:       game.Title,
		Platform:    game.Platform,
		Genre:       game.Genre,
		Alias:       game.Alias,
		Developer:   game.Developer,
		Publisher:   game.Publisher,
		PublishDate: game.PublishDate,
		Thumbnail:   game.Thumbnail,
	}
}

type GameVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Platform    string `json:"platform"`
	Genre       string `json:"genre"`
	Alias       string `json:"alias"`
	Developer   string `json:"developer"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publish_date"`
	Thumbnail   string `json:"thumbnail"`
}
