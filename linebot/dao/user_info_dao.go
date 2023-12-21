package dao

import (
	"linebot/entity"
	"linebot/logger"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
LineのIDを元に、ユーザレコードを取得
*/
func GetUserByLineId(lineId string) *entity.UserInfo {
	db, err := getConnection()

	if err != nil {
		logger.FatalWithStackTrace(err, &logger.LBFT030001)
		return nil
	}

	var ret entity.UserInfo

	db.Table("user_info").Where("line_id = ?", lineId).Find(&ret)
	return &ret
}

/*
LineのIDを元に、最終アクセス時間を更新
UI-A-01
*/
func UpdateUserLastAccess(lineId string) bool {
	db, err := getConnection()

	if err != nil {
		logger.FatalWithStackTrace(err, &logger.LBFT030001)
		return false
	}

	db.Table("user_info").Where("line_id = ?", lineId).Update("last_access", time.Now())
	return true
}

/*
LineのIDを元に、不正なユーザレコードを作成 or 最終アクセス時間を更新
UI-E-01
*/
func UpsertInvalidUser(lineId string) bool {
	db, err := getConnection()

	if err != nil {
		logger.FatalWithStackTrace(err, &logger.LBFT030001)
		return false
	}

	tx := db.Table("user_info")

	tx.Begin()
	var ret entity.UserInfo
	tx.Where("line_id = ?", lineId).Find(&ret)
	if ret.UserUuid == "" {
		ret.UserUuid = uuid.New().String()
	}
	ret.LineId = lineId
	if ret.UserName == "" {
		ret.UserName = "Unknown"
	}
	ret.LastAccess = time.Now()
	tx.Save(ret)
	tx.Commit()
	return true
}

func getConnection() (*gorm.DB, error) {
	database_host := os.Getenv("DATABASE_HOST")
	dsn := "root:mysql@tcp(" + database_host + ":3306)/smart_key?parseTime=true"

	return gorm.Open(mysql.Open(dsn))
}
