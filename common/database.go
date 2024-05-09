package common

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}
	//DB.AutoMigrate(&model.User{})
	//DB.AutoMigrate(&model.Video{})
	//DB.AutoMigrate(&model.UserVideo{})
	//DB.AutoMigrate(&model.Favorite{})
	//DB.AutoMigrate(&model.Comment{})
	//DB.AutoMigrate(&model.Follow{})
	DB.AutoMigrate(&model.Message{})
	fmt.Println("mysql inited")
}
