package dao

import (
  "fmt"
  "encoding/json"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "smart_back/entity"
)

func GetAll(){
  dsn := "root:mysql@tcp(127.0.0.1:3306)/smart_key?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn),&gorm.Config{})
  if err != nil {
    println(err)
  }
  var operationHistory []entity.OperationHistory
  db.Table("operation_history").Find(&operationHistory)
  for i,v := range operationHistory {
    jsonData,err := json.Marshal(v)
    if err != nil{
     fmt.Println(err)
    }
    fmt.Println(i)
    fmt.Printf("%s\n",jsonData)
  }
}
