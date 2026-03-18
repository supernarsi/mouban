package model

import (
	"time"
)

type Song struct {
	ID          uint64 `gorm:"comment:自增主键 ID"`
	DoubanId    uint64 `gorm:"not null;uniqueIndex;comment:豆瓣音乐 ID"`
	Title       string `gorm:"not null;type:varchar(512);comment:专辑名称"`
	Alias       string `gorm:"type:varchar(512);comment:又名"`
	Musician    string `gorm:"type:varchar(2048);comment:音乐人/乐队"`
	AlbumType   string `gorm:"type:varchar(512);comment:专辑类型 (录音室/现场/精选等)"`
	Genre       string `gorm:"type:varchar(512);comment:流派"`
	Media       string `gorm:"type:varchar(512);comment:介质 (CD/黑胶/数字等)"`
	Barcode     string `gorm:"type:varchar(512);comment:条形码"`
	Publisher   string `gorm:"type:varchar(512);comment:出版者"`
	PublishDate string `gorm:"type:varchar(512);comment:发行时间"`
	ISRC        string `gorm:"type:varchar(512);comment:ISRC 编码"`
	AlbumCount  uint32 `gorm:"comment:唱片数量"`
	Intro       string `gorm:"type:mediumtext;comment:专辑简介"`
	TrackList   string `gorm:"type:mediumtext;comment:曲目列表"`
	Thumbnail   string `gorm:"type:varchar(512);comment:封面图 URL"`
	CreatedAt   time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt   time.Time `gorm:"comment:记录更新时间"`
}

func (Song) TableName() string {
	return "song"
}

func (song Song) Show() *SongVO {
	return &SongVO{
		DoubanId:    song.DoubanId,
		Title:       song.Title,
		Alias:       song.Alias,
		Musician:    song.Musician,
		Thumbnail:   song.Thumbnail,
		AlbumType:   song.AlbumType,
		Genre:       song.Genre,
		Media:       song.Media,
		Publisher:   song.Publisher,
		PublishDate: song.PublishDate,
	}
}

type SongVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Alias       string `json:"alias"`
	Musician    string `json:"musician"`
	Thumbnail   string `json:"thumbnail"`
	AlbumType   string `json:"album_type"`
	Genre       string `json:"genre"`
	Media       string `json:"media"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publish_date"`
}
