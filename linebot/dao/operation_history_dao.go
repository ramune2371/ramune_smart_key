package dao

import (
	"fmt"
	"linebot/entity"
	"linebot/logger"
)

// Operation HistoryのInsert
// args: lineId,operationType,operationResult
// return: insertしたレコード
func (g *GormDB) InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) *entity.OperationHistory {
	table := g.getTable(entity.OperationHistoryTable)
	if table == nil {
		return nil
	}

	data := entity.OperationHistory{
		LineId:          lineId,
		OperationType:   operationType,
		OperationResult: operationResult,
	}
	table.Create(&data)
	table.First(&data)
	logger.Debug(fmt.Sprintf("data: %+v", data))
	return &data
}

// Operation HistoryのOperation Resultカラムの更新
// エラー時はUpdateOperationHistoryWithErrorByOperationIdを使用すること
// args: operationId=更新対象のoperationId, result=更新するOperationResult
// return: 更新したレコードのID(error時は-1)
func (g *GormDB) UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) int {
	table := g.getTable(entity.OperationHistoryTable)

	if table == nil {
		return -1
	}
	var target entity.OperationHistory
	table.Where("operation_id = ?", operationId).First(&target)
	target.OperationResult = result
	table.Save(&target)
	return *target.OperationId
}

// Operation HistoryのerrorCodeカラムの更新
// OperationResultカラムにはentity.Errorが入る
// args: operationId=更新対象のOperationId, errorCode=更新するOperationErrorCode
// return 更新したOperationId。エラー時は-1
func (g *GormDB) UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) int {
	table := g.getTable(entity.OperationHistoryTable)
	if table == nil {
		return -1
	}
	var target entity.OperationHistory
	table.Where("operation_id = ?", operationId).Find(&target)
	target.ErrorCode = errorCode
	target.OperationResult = entity.Error
	table.Save(target)
	return *target.OperationId
}
