package dao

import (
	"mouban/common"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func CountBook() int64 {
	var count int64
	result := common.Db.Model(&model.Book{}).Count(&count)
	if result.Error != nil {
		logrus.Errorln("count book error:", result.Error)
		return 0
	}
	return count
}

func UpsertBook(book *model.Book) {
	logrus.WithField("upsert", "book").Infoln("upsert book", book.DoubanId, book.Title)
	data := &model.Book{}
	result := common.Db.Where("douban_id = ?", book.DoubanId).Assign(book).FirstOrCreate(data)
	if result.Error != nil {
		logrus.Errorln("upsert book error:", result.Error, "douban_id:", book.DoubanId)
	}
}

func UpdateBookThumbnail(doubanId uint64, thumbnail string) {
	result := common.Db.Model(&model.Book{}).Where("douban_id = ?", doubanId).Update("thumbnail", thumbnail)
	if result.Error != nil {
		logrus.Errorln("update book thumbnail error:", result.Error, "douban_id:", doubanId)
	}
}

func CreateBookNx(book *model.Book) bool {
	data := &model.Book{}
	inserted := common.Db.Where("douban_id = ?", book.DoubanId).Attrs(book).FirstOrCreate(data).RowsAffected > 0
	if inserted {
		logrus.Infoln("create book", book.DoubanId, book.Title)
	}
	return inserted
}

func GetBookDetail(doubanId uint64) *model.Book {
	book := &model.Book{}
	result := common.Db.Where("douban_id = ?", doubanId).Find(book)
	if result.Error != nil {
		logrus.Errorln("get book detail error:", result.Error, "douban_id:", doubanId)
		return nil
	}
	if book.ID == 0 {
		return nil
	}
	return book
}

func ListBookBrief(doubanIds *[]uint64) *[]model.Book {
	var books *[]model.Book
	result := common.Db.Omit("serial", "isbn", "framing", "page", "intro").Where("douban_id IN ?", *doubanIds).Find(&books)
	if result.Error != nil {
		logrus.Errorln("list book brief error:", result.Error)
		return nil
	}
	return books
}
