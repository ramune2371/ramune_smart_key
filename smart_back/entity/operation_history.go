package entity

import "time"

type OperationHistory struct{
  OperationId int
  CertificateUuid string
  OperationType int
  OperationResult bool
  ErrorCode string
  OperationTime time.Time
}
