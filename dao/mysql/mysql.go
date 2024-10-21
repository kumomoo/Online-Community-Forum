package mysql

import (
	"bluebell/models"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect DB failed: %v", zap.Error(err))
		return
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Community{})
	db.AutoMigrate(&models.Post{})
	fmt.Printf("connect mysql success!\n")
	return
}

func Close() {
	DB, err := db.DB()
	if err != nil {
		zap.L().Fatal("close mysql failed: %v", zap.Error(err))
	}
	DB.Close()
}
