package entity

import (
	"fmt"
)

type OperationType int

const (
	Open  OperationType = 0
	Close OperationType = 1
	Check OperationType = 2
)

type Operation struct {
	OperationId int
	UserId      string
	Operation   OperationType
	ReplyToken  string
}

func (op Operation) String() string {
	return fmt.Sprintf("OperationId: %d, UserId:%s, Operation:%+v, ReplyToken:%s", op.OperationId, op.UserId, op.Operation, op.ReplyToken)
}

func (op Operation) IsEqual(target Operation) bool {
	return op.Operation == target.Operation &&
		op.OperationId == target.OperationId &&
		op.ReplyToken == target.ReplyToken &&
		op.UserId == target.UserId
}
