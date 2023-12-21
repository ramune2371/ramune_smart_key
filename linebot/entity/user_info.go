package entity

import "time"

type UserInfo struct {
	UserUuid   string
	LineId     string
	UserName   string
	LastAccess time.Time `gorm:"autoCreateTime:true"`
	Active     bool
}
