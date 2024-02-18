package user_info

import (
	"linebot/dao/database"
	"linebot/entity"
	"linebot/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserInfoDaoImpl struct {
	Database database.DatabaseConnection
}

/*
LineのIDを元に、ユーザレコードを取得
*/
func (uiDao UserInfoDaoImpl) GetUserByLineId(lineId string) *entity.UserInfo {
	var ret *entity.UserInfo

	res := uiDao.Database.ReadOnly(entity.UserInfoTable, func(tx *gorm.DB) error {
		if err := tx.Where("line_id = ?", lineId).Find(&ret).Error; err != nil {
			logger.ErrorWithStackTrace(err, logger.LBER030001)
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
func (uiDao UserInfoDaoImpl) UpdateUserLastAccess(lineId string) bool {

	res := uiDao.Database.ReadWrite(entity.UserInfoTable, func(tx *gorm.DB) error {
		if err := tx.Where("line_id = ?", lineId).Update("last_access", time.Now()).Error; err != nil {
			logger.ErrorWithStackTrace(err, logger.LBER030002)
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
func (uiDao UserInfoDaoImpl) UpsertInvalidUser(lineId string) bool {

	var ret *entity.UserInfo
	res := uiDao.Database.ReadWrite(entity.UserInfoTable, func(tx *gorm.DB) error {
		if err := tx.Where("line_id = ?", lineId).Find(&ret).Error; err != nil {
			logger.ErrorWithStackTrace(err, logger.LBER030001)
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
		if err := tx.Where("line_id = ?", lineId).Save(&ret).Error; err != nil {
			logger.ErrorWithStackTrace(err, logger.LBER030002)
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
