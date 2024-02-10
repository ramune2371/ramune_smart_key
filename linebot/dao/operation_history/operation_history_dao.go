package operation_history

import (
	"linebot/entity"
)

type OperationHistoryDao interface {
	InsertOperationHistory(lineId string, operationType entity.OperationType, operationResult entity.OperationResult) *entity.OperationHistory
	UpdateOperationHistoryByOperationId(operationId int, result entity.OperationResult) int
	UpdateOperationHistoryWithErrorByOperationId(operationId int, errorCode entity.OperationErrorCode) int
}
