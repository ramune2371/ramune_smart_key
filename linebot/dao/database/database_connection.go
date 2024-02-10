package database

import "gorm.io/gorm"

type DatabaseConnection interface {
	ReadWrite(tableName string, fc func(tx *gorm.DB) error) error
	ReadOnly(tableName string, fc func(tx *gorm.DB) error) error
}
