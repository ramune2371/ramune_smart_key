package dao

import (
	"linebot/entity"
	"linebot/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetUserByLineId(lineId string) *entity.UserInfo {
  dsn := "root:mysql@tcp(127.0.0.1:3306)/smart_key"
  db,err := gorm.Open(mysql.Open(dsn))

  if err != nil {
    logger.Fatal("DBとの接続に失敗しました。","LBFT03001")
    return nil
  }

  var ret entity.UserInfo

  db.Table("user_info").Where("line_id = ?",lineId).Find(&ret)
  return &ret
}
