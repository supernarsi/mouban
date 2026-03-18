package model

import (
	"time"
)

type User struct {
	ID           uint64 `gorm:"comment:自增主键 ID"`
	DoubanUid    uint64 `gorm:"not null;uniqueIndex;comment:豆瓣用户 ID"`
	Domain       string `gorm:"not null;index;type:varchar(64);comment:豆瓣个人主页域名 (如 ahbei)"`
	Name         string `gorm:"not null;type:varchar(512);comment:用户昵称"`
	Thumbnail    string `gorm:"type:varchar(512);comment:用户头像 URL"`
	BookWish     uint32 `gorm:"not null default 0;comment:想读数量"`
	BookDo       uint32 `gorm:"not null default 0;comment:在读数量"`
	BookCollect  uint32 `gorm:"not null default 0;comment:读过数量"`
	GameWish     uint32 `gorm:"not null default 0;comment:想玩数量"`
	GameDo       uint32 `gorm:"not null default 0;comment:在玩数量"`
	GameCollect  uint32 `gorm:"not null default 0;comment:玩过数量"`
	MovieWish    uint32 `gorm:"not null default 0;comment:想看数量"`
	MovieDo      uint32 `gorm:"not null default 0;comment:在看数量"`
	MovieCollect uint32 `gorm:"not null default 0;comment:看过数量"`
	SongWish     uint32 `gorm:"not null default 0;comment:想听数量"`
	SongDo       uint32 `gorm:"not null default 0;comment:在听数量"`
	SongCollect  uint32 `gorm:"not null default 0;comment:听过数量"`
	SyncAt       time.Time `gorm:"comment:最近同步时间"`
	CheckAt      time.Time `gorm:"comment:最近检测时间"`
	RegisterAt   time.Time `gorm:"comment:注册时间 (首次抓取时间)"`
	PublishAt    time.Time `gorm:"comment:用户最近发布时间 (用于判断是否变化)"`
	CreatedAt    time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt    time.Time `gorm:"comment:记录更新时间"`
}

func (User) TableName() string {
	return "user"
}

func (user User) Show() *UserVO {
	return &UserVO{
		ID:           user.DoubanUid,
		Domain:       user.Domain,
		Name:         user.Name,
		Thumbnail:    user.Thumbnail,
		BookWish:     user.BookWish,
		BookDo:       user.BookDo,
		BookCollect:  user.BookCollect,
		GameWish:     user.GameWish,
		GameDo:       user.GameDo,
		GameCollect:  user.GameCollect,
		MovieWish:    user.MovieWish,
		MovieDo:      user.MovieDo,
		MovieCollect: user.MovieCollect,
		SongWish:     user.SongWish,
		SongDo:       user.SongDo,
		SongCollect:  user.SongCollect,
		PublishAt:    user.PublishAt.Unix(),
		SyncAt:       user.SyncAt.Unix(),
		CheckAt:      user.CheckAt.Unix(),
	}
}

type UserVO struct {
	ID           uint64 `json:"id"`
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	Thumbnail    string `json:"thumbnail"`
	BookWish     uint32 `json:"book_wish"`
	BookDo       uint32 `json:"book_do"`
	BookCollect  uint32 `json:"book_collect"`
	GameWish     uint32 `json:"game_wish"`
	GameDo       uint32 `json:"game_do"`
	GameCollect  uint32 `json:"game_collect"`
	MovieWish    uint32 `json:"movie_wish"`
	MovieDo      uint32 `json:"movie_do"`
	MovieCollect uint32 `json:"movie_collect"`
	SongWish     uint32 `json:"song_wish"`
	SongDo       uint32 `json:"song_do"`
	SongCollect  uint32 `json:"song_collect"`
	PublishAt    int64  `json:"publish_at"`
	SyncAt       int64  `json:"sync_at"`
	CheckAt      int64  `json:"check_at"`
}
