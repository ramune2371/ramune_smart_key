package dao

import (
  "fmt"
	"linebot/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserInfo struct {
  UserUuuid string
  LineId string
  UserName string
  LastAccess string
  Active bool
}

func GetUserActive(userId string) (string,bool) {
  dsn := "root:mysql@tcp(127.0.0.1:3306)/smart_key"
  db,err := gorm.Open(mysql.Open(dsn))
  if err != nil {
    logger.Fatal("DBとの接続に失敗しました。","LBFT03001")
    return "",false
  }

  var rets []UserInfo

  db.Table("user_info").Where("line_id = ?",userId).Find(&rets)
  for _,r := range rets {
    logger.Debug(fmt.Sprintf("user name = %s",r.UserName))
  }
  if len(rets) != 1 {
    return "",false
  }
  return rets[0].UserName,rets[0].Active
}
