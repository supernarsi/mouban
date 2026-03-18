package model

import (
	"time"
)

type Schedule struct {
	ID        uint64    `gorm:"comment:自增主键 ID"`
	DoubanId  uint64    `gorm:"not null;uniqueIndex:uk_schedule;comment:目标 ID (用户 ID 或条目 ID)"`
	Type      uint8     `gorm:"not null;uniqueIndex:uk_schedule;index:idx_status;index:idx_result;index:idx_search;priority=1;comment:类型：0=用户，1=书，2=电影，3=游戏，4=音乐"`
	Status    *uint8    `gorm:"not null;index:idx_status;index:idx_search;priority:2;comment:爬取状态：0=待爬取，1=爬取中，2=已爬取，3=可爬取"`
	Result    *uint8    `gorm:"not null;index:idx_result;index:idx_search;priority:2;comment:爬取结果：0=未就绪，1=就绪，2=无效"`
	CreatedAt time.Time `gorm:"comment:任务创建时间"`
	UpdatedAt time.Time `gorm:"index:idx_result;index:idx_status;index:idx_search;priority:3;comment:任务更新时间"`
}

func (Schedule) TableName() string {
	return "schedule"
}
