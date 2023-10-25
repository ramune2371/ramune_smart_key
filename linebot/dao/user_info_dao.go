package dao

import (
	"linebot/entity"
	"linebot/logger"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetUserByLineId(lineId string) *entity.UserInfo {
	database_host := os.Getenv("DATABASE_HOST")
	dsn := "root:mysql@tcp(" + database_host + ":3306)/smart_key"
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		logger.FatalWithStackTrace(err, &logger.LBFT030001)
		return nil
	}

	var ret entity.UserInfo

	db.Table("user_info").Where("line_id = ?", lineId).Find(&ret)
	return &ret
}
