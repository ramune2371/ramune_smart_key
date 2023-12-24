package dao

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDB struct {
	DB *gorm.DB
}

func InitDB() *GormDB {
	database_host := os.Getenv("DATABASE_HOST")
	dsn := "root:mysql@tcp(" + database_host + ":3306)/smart_key?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("DB Connection Error")
	}
	return &GormDB{DB: db}
}

func (g *GormDB) getTable(tableName string) *gorm.DB {

	return g.DB.Table(tableName)
}
