package entity

type OperationResult int

const (
	Error     OperationResult = -1
	Operating OperationResult = 0
	Success   OperationResult = 1
	Merged    OperationResult = 2
)
