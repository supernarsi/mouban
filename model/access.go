package model

import (
	"time"
)

type Access struct {
	ID        uint64 `gorm:"comment:自增主键 ID"`
	DoubanUid uint64 `gorm:"not null;index;comment:豆瓣用户 ID (请求参数)"`
	Path      string `gorm:"not null;type:varchar(64);comment:请求路径"`
	Ip        string `gorm:"not null;type:varchar(64);comment:请求 IP"`
	UserAgent string `gorm:"not null;type:varchar(512);comment:User-Agent"`
	Referer   string `gorm:"not null;type:varchar(512);comment:Referer 来源"`
	CreatedAt time.Time `gorm:"comment:访问时间"`
	UpdatedAt time.Time `gorm:"comment:记录更新时间"`
}

func (Access) TableName() string {
	return "access"
}
