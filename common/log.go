package common

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitLogger() {
	// 从配置读取日志级别
	levelStr := strings.ToLower(viper.GetString("log.level"))
	var level logrus.Level
	switch levelStr {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn", "warning":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	default:
		level = logrus.InfoLevel
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
		TimestampFormat:   "2006-01-02 15:04:05",
	})

	logrus.WithField("level", level.String()).Infoln("logrus init success")
}
