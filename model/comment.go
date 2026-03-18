package model

import (
	"time"
)

type Comment struct {
	ID        uint64 `gorm:"comment:自增主键 ID"`
	DoubanUid uint64 `gorm:"not null;uniqueIndex:uk_comment;index:idx_search;priority:1;comment:豆瓣用户 ID"`
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_comment;priority:2;comment:条目 ID (书/影/音/游)"`
	Type      uint8  `gorm:"not null;uniqueIndex:uk_comment;index:idx_search;priority:3;comment:条目类型：0=用户，1=书，2=电影，3=游戏，4=音乐"`
	Rate      uint8  `gorm:"comment:评分 (0-5 星，0 表示未评分)"`
	Label     string `gorm:"type:varchar(512);comment:标签/关键词"`
	Comment   string `gorm:"type:mediumtext;comment:评论内容"`
	Action    *uint8 `gorm:"not null;index:idx_search;priority:4;comment:操作类型：0=do,1=wish,2=collect,3=hide"`
	MarkDate  time.Time `gorm:"not null;index:idx_search;priority:5;comment:用户标记时间 (评论发布日期)"`
	CreatedAt time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt time.Time `gorm:"comment:记录更新时间"`
}

func (Comment) TableName() string {
	return "comment"
}

func (comment Comment) Show(item interface{}) *CommentVO {
	return &CommentVO{
		Item:     item,
		Rate:     comment.Rate,
		Label:    comment.Label,
		Comment:  comment.Comment,
		Action:   *comment.Action,
		MarkDate: comment.MarkDate.Format("2006-01-02"),
	}
}

type CommentVO struct {
	Item     interface{} `json:"item"`
	Rate     uint8       `json:"rate"`
	Label    string      `json:"label"`
	Comment  string      `json:"comment"`
	Action   uint8       `json:"action"`
	MarkDate string      `json:"mark_date"`
}
