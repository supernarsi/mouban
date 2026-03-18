package common

import (
	"database/sql"
	"fmt"
	"mouban/model"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB
var sqlDB *sql.DB

func InitDatabase() {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")

	tryCreateDB(username, password, host, port, database)
	getConnection(username, password, host, port, database, charset, loc)
	migrateTables()
}

func tryCreateDB(username string, password string, host string, port string, database string) {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)

	db, err := sql.Open("mysql", sqlStr)
	if err != nil {
		logrus.Errorln("open database failed:", err)
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logrus.Infoln("database close failed")
		}
	}(db)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s ;", database))
	if err != nil {
		logrus.Errorln("create database failed:", err)
		panic(err)
	}
}

func getConnection(username string, password string, host string, port string, database string, charset string, loc string) {

	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))

	dbLogger := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             500 * time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	gormDB, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		logrus.Infoln("Open database failed", err)
		panic("Open database failed " + err.Error())
	}
	Db = gormDB

	// 获取底层 sql.DB 并配置连接池
	sqlDB, err = gormDB.DB()
	if err != nil {
		logrus.Infoln("Get underlying DB failed", err)
		panic("Get underlying DB failed " + err.Error())
	}

	// 配置连接池参数
	sqlDB.SetMaxIdleConns(10)              // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)             // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)    // 连接最长生命周期

	logrus.Infoln("mysql connect success, connection pool configured")
}

func migrateTables() {
	migrator := Db.Migrator()

	// 表注释定义
	tableComments := map[string]string{
		"access":   "访问日志表 - 记录 API 访问日志，用于限流和审计",
		"book":     "书籍条目信息表 - 存储书籍详细信息",
		"comment":  "用户评论表 - 存储用户对条目的评论/标注/评分记录",
		"game":     "游戏条目信息表 - 存储游戏详细信息",
		"movie":    "电影条目信息表 - 存储电影详细信息",
		"song":     "音乐专辑信息表 - 存储音乐专辑详细信息",
		"rating":   "评分统计表 - 存储条目的聚合评分数据",
		"schedule": "爬虫调度队列表 - 存储待爬取任务队列",
		"user":     "豆瓣用户信息表 - 存储用户基本信息和各类目数量统计",
		"storage":  "存储映射表 - 记录原始图片 URL 到 S3 存储的映射关系",
	}

	// 只在表不存在时创建表，避免覆盖已有表结构
	tables := []struct {
		model interface{ TableName() string }
		name  string
	}{
		{&model.Access{}, "access"},
		{&model.Book{}, "book"},
		{&model.Comment{}, "comment"},
		{&model.Game{}, "game"},
		{&model.Movie{}, "movie"},
		{&model.Song{}, "song"},
		{&model.Rating{}, "rating"},
		{&model.Schedule{}, "schedule"},
		{&model.User{}, "user"},
		{&model.Storage{}, "storage"},
	}

	for _, t := range tables {
		if !migrator.HasTable(t.model) {
			// 创建表并设置表注释
			if err := migrator.CreateTable(t.model); err != nil {
				panic("create table " + t.name + " failed: " + err.Error())
			}
			// 为表添加注释
			comment, ok := tableComments[t.name]
			if ok {
				if err := Db.Exec(fmt.Sprintf("ALTER TABLE %s COMMENT = ?", t.name), comment).Error; err != nil {
					logrus.Warnln("add comment to table " + t.name + " failed:", err)
				}
			}
			logrus.Infoln("create table " + t.name + " success")
		} else {
			logrus.Infoln("table " + t.name + " already exists, skip creation")
		}
	}
}
