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

	err := Db.AutoMigrate(
		&model.Access{},
		&model.Book{},
		&model.Comment{},
		&model.Game{},
		&model.Movie{},
		&model.Song{},
		&model.Rating{},
		&model.Schedule{},
		&model.User{},
		&model.Storage{},
	)
	if err != nil {
		panic("init database failed " + err.Error())
	}

}
