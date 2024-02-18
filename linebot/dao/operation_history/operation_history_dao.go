package operation_history

import (
	"linebot/entity"
)

type OperationHistoryDao interface {
	InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) (*entity.OperationHistory, error)
	UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) (int, error)
	UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) (int, error)
}

type emptyOperationHistoryDao struct{}

func (emptyOperationHistoryDao) InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) (*entity.OperationHistory, error) {
	return nil, nil
}
func (emptyOperationHistoryDao) UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) (int, error) {
	return -1, nil
}
func (emptyOperationHistoryDao) UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) (int, error) {
	return -1, nil
}

func NewEmptyOperationHistoryDao() *emptyOperationHistoryDao {
	return &emptyOperationHistoryDao{}
}
