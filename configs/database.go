package configs

import (
	// "fmt"

	"echo-api/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func RunDatabase(dbName string) (*gorm.DB, error) {
	var localLog = logrus.New()
	err := godotenv.Load()
	if err != nil {
		localLog.Warn("Error loading .env file, using default environment variables")
	}

	cfg := &models.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName,
	)

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		localLog.Error("Failed to connect to database:", err)
		return nil, err
	}

	localLog.Infof("Database '%s' connected successfully", dbName)
	return db, nil
}
