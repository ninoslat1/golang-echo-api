package services

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	log := logrus.New()
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", 0755)
		if err != nil {
			log.Error("Failed to create log directory: ERROR" + err.Error())
		} else {
			log.Info("Log directory created")
		}
	}

	logSetup := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, logSetup)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	log.SetOutput(multiWriter)

	log.SetLevel(logrus.InfoLevel)

	log.Info("Logger initialized successfully")

	return log
}
