package operation_history

import (
	"fmt"
	"linebot/applicationerror"
	"linebot/dao/database"
	"linebot/entity"
	"linebot/logger"

	"gorm.io/gorm"
)

type OperationHistoryDaoImpl struct {
	Database database.DatabaseConnection
}

/*
# Operation HistoryのInsert

args: lineId,operationType,operationResult

return: insertしたレコード
*/
func (OperationHistoryDaoImpl OperationHistoryDaoImpl) InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) (*entity.OperationHistory, error) {

	data := entity.OperationHistory{
		LineId:          lineId,
		OperationType:   operationType,
		OperationResult: operationResult,
	}

	err := OperationHistoryDaoImpl.Database.ReadWrite(entity.OperationHistoryTable, func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			logger.ErrorWithStackTrace(err, applicationerror.DBInsertError, logger.LBER030002)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Debug(fmt.Sprintf("data: %+v", data))
	return &data, nil
}

/*
# Operation HistoryのOperation Resultカラムの更新

エラー時はUpdateOperationHistoryWithErrorByOperationIdを使用すること

args: operationId=更新対象のoperationId, result=更新するOperationResult

return: 更新したレコードのID(error時は-1)
*/
func (OperationHistoryDaoImpl OperationHistoryDaoImpl) UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) (int, error) {

	var target entity.OperationHistory

	err := OperationHistoryDaoImpl.Database.ReadWrite(entity.OperationHistoryTable, func(tx *gorm.DB) error {
		if err := tx.Where("operation_id = ?", operationId).First(&target).Error; err != nil {
			logger.ErrorWithStackTrace(err, applicationerror.DBSelectError, logger.LBER030001)
			return err
		}
		target.OperationResult = result
		if err := tx.Save(&target).Error; err != nil {
			logger.ErrorWithStackTrace(err, applicationerror.DBUpdateError, logger.LBER030002)
			return err
		}
		return nil
	})

	if err != nil || target.OperationId == nil {
		return -1, err
	}

	return *target.OperationId, nil
}

/*
# Operation HistoryのerrorCodeカラムの更新

# OperationResultカラムにはentity.Errorが入る

args: operationId=更新対象のOperationId, errorCode=更新するOperationErrorCode

return 更新したOperationId。エラー時は-1
*/
func (OperationHistoryDaoImpl OperationHistoryDaoImpl) UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) (int, error) {
	var target entity.OperationHistory
	err := OperationHistoryDaoImpl.Database.ReadWrite(entity.OperationHistoryTable, func(tx *gorm.DB) error {
		if err := tx.Where("operation_id = ?", operationId).Find(&target).Error; err != nil {
			logger.ErrorWithStackTrace(err, applicationerror.DBSelectError, logger.LBER030001)
			return err
		}
		target.ErrorCode = errorCode
		target.OperationResult = entity.Error
		if err := tx.Save(&target).Error; err != nil {
			logger.ErrorWithStackTrace(err, applicationerror.DBUpdateError, logger.LBER030001)
			return err
		}
		return nil
	})
	if err == nil {
		return -1, err
	}
	return *target.OperationId, nil
}
