package entity

import "time"

type OperationErrorCode string

const (
	OperatingTypeError       OperationErrorCode = "302"
	KeyServerConnectionError OperationErrorCode = "303"
	InOperatingError         OperationErrorCode = "304"
)

type OperationHistory struct {
	OperationId     *int `gorm:"primaryKey"`
	LineId          string
	OperationType   OperationType
	OperationResult OperationResult
	ErrorCode       OperationErrorCode
	OperationTime   *time.Time `gorm:"autoCreateTime"`
}
