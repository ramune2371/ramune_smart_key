package dao

import (
	"linebot/entity"
	"linebot/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
LineのIDを元に、ユーザレコードを取得
*/
func GetUserByLineId(lineId string) *entity.UserInfo {
	var ret *entity.UserInfo

	res := readOnly(entity.UserInfoTable, func(tx *gorm.DB) error {
		if err := tx.Where("line_id = ?", lineId).Find(&ret).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030001)
			return err
		}
		return nil
	})

	if res != nil {
		return nil
	}

	return ret
}

/*
LineのIDを元に、最終アクセス時間を更新
UI-A-01
*/
func UpdateUserLastAccess(lineId string) bool {

	res := readWrite(entity.UserInfoTable, func(tx *gorm.DB) error {
		if err := tx.Where("line_id = ?", lineId).Update("last_access", time.Now()).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030002)
			return err
		}
		return nil
	})

	if res != nil {
		return false
	} else {
		return true
	}
}

/*
LineのIDを元に、不正なユーザレコードを作成 or 最終アクセス時間を更新
UI-E-01
*/
func UpsertInvalidUser(lineId string) bool {

	var ret *entity.UserInfo
	res := readWrite(entity.UserInfoTable, func(tx *gorm.DB) error {
		if err := tx.Where("line_id = ?", lineId).Find(&ret).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030001)
			return err
		}
		if ret.UserUuid == "" {
			ret.UserUuid = uuid.New().String()
		}
		ret.LineId = lineId
		if ret.UserName == "" {
			ret.UserName = "Unknown"
		}
		ret.Active = false
		if err := tx.Save(&ret).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030002)
			return err
		}
		return nil
	})
	if res != nil {
		return false
	} else {
		return true
	}
}
