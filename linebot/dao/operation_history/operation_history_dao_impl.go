package operation_history

import (
	"fmt"
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
func (OperationHistoryDaoImpl OperationHistoryDaoImpl) InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) *entity.OperationHistory {

	data := entity.OperationHistory{
		LineId:          lineId,
		OperationType:   operationType,
		OperationResult: operationResult,
	}

	res := OperationHistoryDaoImpl.Database.ReadWrite(entity.OperationHistoryTable, func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030002)
			return err
		}
		return nil
	})

	if res != nil {
		return nil
	}

	logger.Debug(fmt.Sprintf("data: %+v", data))
	return &data
}

/*
# Operation HistoryのOperation Resultカラムの更新

エラー時はUpdateOperationHistoryWithErrorByOperationIdを使用すること

args: operationId=更新対象のoperationId, result=更新するOperationResult

return: 更新したレコードのID(error時は-1)
*/
func (OperationHistoryDaoImpl OperationHistoryDaoImpl) UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) int {

	var target entity.OperationHistory

	res := OperationHistoryDaoImpl.Database.ReadWrite(entity.OperationHistoryTable, func(tx *gorm.DB) error {
		if err := tx.Where("operation_id = ?", operationId).First(&target).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030001)
			return err
		}
		target.OperationResult = result
		if err := tx.Save(&target).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030002)
			return err
		}
		return nil
	})

	if res != nil {
		return -1
	}

	return *target.OperationId
}

/*
# Operation HistoryのerrorCodeカラムの更新

# OperationResultカラムにはentity.Errorが入る

args: operationId=更新対象のOperationId, errorCode=更新するOperationErrorCode

return 更新したOperationId。エラー時は-1
*/
func (OperationHistoryDaoImpl OperationHistoryDaoImpl) UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) int {
	var target entity.OperationHistory
	res := OperationHistoryDaoImpl.Database.ReadWrite(entity.OperationHistoryTable, func(tx *gorm.DB) error {
		if err := tx.Where("operation_id = ?", operationId).Find(&target).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030001)
			return err
		}
		target.ErrorCode = errorCode
		target.OperationResult = entity.Error
		if err := tx.Save(target).Error; err != nil {
			logger.ErrorWithStackTrace(err, &logger.LBER030001)
			return err
		}
		return nil
	})
	if res == nil {
		return -1
	}
	return *target.OperationId
}
