package dao

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	database_host := os.Getenv("DATABASE_HOST")
	dsn := "root:mysql@tcp(" + database_host + ":3306)/smart_key?parseTime=true"
	database, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	db = database
}

func getTable(tableName string) *gorm.DB {

	return db.Table(tableName)
}

func Close() {
	sqldb, _ := db.DB()
	sqldb.Close()
}
