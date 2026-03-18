package model

import (
	"time"
)

type Storage struct {
	ID        uint64 `gorm:"comment:自增主键 ID"`
	Source    string `gorm:"type:varchar(256);not null;uniqueIndex;comment:原始图片 URL"`
	Target    string `gorm:"type:varchar(256);not null;comment:S3 存储后的 URL"`
	Md5       string `gorm:"type:varchar(64);not null;index;comment:图片内容 MD5 值 (用于去重)"`
	CreatedAt time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt time.Time `gorm:"comment:记录更新时间"`
}

func (Storage) TableName() string {
	return "storage"
}
