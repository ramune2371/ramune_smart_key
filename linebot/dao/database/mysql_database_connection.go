package database

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type MySQLDatabaseConnection struct {
	db *gorm.DB
}

func (dbConnection *MySQLDatabaseConnection) InitDB() {
	database_host := os.Getenv("DATABASE_HOST")
	dsn := "root:mysql@tcp(" + database_host + ":3306)/smart_key?parseTime=true"
	database, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	dbConnection.db = database
}

func (dbConnection MySQLDatabaseConnection) getTable(tableName string) *gorm.DB {

	return dbConnection.db.Table(tableName)
}

func (dbConnection MySQLDatabaseConnection) ReadWrite(tableName string, fc func(tx *gorm.DB) error) error {
	return dbConnection.getTable(tableName).Clauses(dbresolver.Write).Transaction(fc)
}

func (dbConnection MySQLDatabaseConnection) ReadOnly(tableName string, fc func(tx *gorm.DB) error) error {
	return dbConnection.getTable(tableName).Clauses(dbresolver.Read).Transaction(fc)
}

func (dbConnection *MySQLDatabaseConnection) Close() {
	sqldb, _ := dbConnection.db.DB()
	sqldb.Close()
}
