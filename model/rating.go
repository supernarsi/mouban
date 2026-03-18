package model

import (
	"time"
)

type Rating struct {
	ID        uint64 `gorm:"comment:自增主键 ID"`
	Type      uint8  `gorm:"not null;uniqueIndex:uk_unique_id;comment:条目类型：1=书，2=电影，3=游戏，4=音乐"`
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_unique_id;comment:条目 ID"`
	Total     uint32 `gorm:"comment:评分总人数"`
	Rating    float32 `gorm:"comment:平均分 (0-10 分)"`
	Star5     float32 `gorm:"comment:5 星占比 (百分比)"`
	Star4     float32 `gorm:"comment:4 星占比 (百分比)"`
	Star3     float32 `gorm:"comment:3 星占比 (百分比)"`
	Star2     float32 `gorm:"comment:2 星占比 (百分比)"`
	Star1     float32 `gorm:"comment:1 星占比 (百分比)"`
	Status    *uint8 `gorm:"comment:状态：0=正常，1=人数不足，2=不允许显示"`
	CreatedAt time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt time.Time `gorm:"comment:记录更新时间"`
}

func (Rating) TableName() string {
	return "rating"
}

type RatingVO struct {
	Total  uint32  `json:"total"`
	Rating float32 `json:"rating"`
	Star5  float32 `json:"star5"`
	Star4  float32 `json:"star4"`
	Star3  float32 `json:"star3"`
	Star2  float32 `json:"star2"`
	Star1  float32 `json:"star1"`
	Status string  `json:"status"`
}
