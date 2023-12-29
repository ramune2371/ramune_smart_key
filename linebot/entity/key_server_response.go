package entity

type KeyStatus string
type OperationStatus string

const (
	KeyStatusOpen  KeyStatus = "True"
	KeyStatusClose KeyStatus = "False"
)

const (
	OperationAnother  OperationStatus = "another"
	OperationComplete OperationStatus = "complete"
	OperationAlready  OperationStatus = "already"
	OperationUnknown  OperationStatus = "unknown"
)

type KeyServerResponse struct {
	KeyStatus       KeyStatus
	OperationStatus OperationStatus
}
