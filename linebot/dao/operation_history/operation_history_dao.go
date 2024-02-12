package operation_history

import (
	"linebot/entity"
)

type OperationHistoryDao interface {
	InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) *entity.OperationHistory
	UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) int
	UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) int
}

type emptyOperationHistoryDao struct{}

func (emptyOperationHistoryDao) InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) *entity.OperationHistory {
	return nil
}
func (emptyOperationHistoryDao) UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) int {
	return -1
}
func (emptyOperationHistoryDao) UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) int {
	return -1
}

func NewEmptyOperationHistoryDao() *emptyOperationHistoryDao {
	return &emptyOperationHistoryDao{}
}
