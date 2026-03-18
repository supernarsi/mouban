package model

import (
	"time"
)

type Book struct {
	ID          uint64 `gorm:"comment:自增主键 ID"`
	DoubanId    uint64 `gorm:"not null;uniqueIndex;comment:豆瓣书籍 ID"`
	Title       string `gorm:"not null;type:varchar(1024);comment:书名"`
	Subtitle    string `gorm:"type:varchar(1024);comment:副标题"`
	Orititle    string `gorm:"type:varchar(1024);comment:原作名 (翻译书籍的原名)"`
	Author      string `gorm:"type:varchar(1024);comment:作者 (多人用分隔符)"`
	Translator  string `gorm:"type:varchar(512);comment:译者"`
	Press       string `gorm:"type:varchar(512);comment:出版社"`
	Producer    string `gorm:"type:varchar(512);comment:出品方"`
	Serial      string `gorm:"type:varchar(512);comment:丛书名"`
	PublishDate string `gorm:"type:varchar(64);comment:出版年月"`
	ISBN        string `gorm:"type:varchar(64);comment:ISBN 号"`
	Framing     string `gorm:"type:varchar(512);comment:装帧 (精装/平装等)"`
	Page        uint32 `gorm:"comment:页数"`
	Price       uint32 `gorm:"comment:定价 (单位：分)"`
	BookIntro   string `gorm:"type:mediumtext;comment:书籍简介"`
	AuthorIntro string `gorm:"type:mediumtext;comment:作者简介"`
	Thumbnail   string `gorm:"type:varchar(512);comment:封面图 URL"`
	CreatedAt   time.Time `gorm:"comment:记录创建时间"`
	UpdatedAt   time.Time `gorm:"comment:记录更新时间"`
}

func (Book) TableName() string {
	return "book"
}

func (book Book) Show() *BookVO {
	return &BookVO{
		DoubanId:    book.DoubanId,
		Title:       book.Title,
		Subtitle:    book.Subtitle,
		Orititle:    book.Orititle,
		Author:      book.Author,
		Translator:  book.Translator,
		Press:       book.Press,
		Producer:    book.Producer,
		PublishDate: book.PublishDate,
		Thumbnail:   book.Thumbnail,
	}
}

type BookVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Orititle    string `json:"orititle"`
	Author      string `json:"author"`
	Translator  string `json:"translator"`
	Press       string `json:"press"`
	Producer    string `json:"producer"`
	PublishDate string `json:"publish_date"`
	Thumbnail   string `json:"thumbnail"`
}
