package applicationerror

// TODO Goのエラー生成うまく理解してない。要勉強＆リファクタ

type ApplicationError struct {
	Code    string
	Message string
}

func (e ApplicationError) Error() string { return e.Message }

func (e ApplicationError) ApplicationError() bool {
	return true
}

var (
	ConnectionError           = ApplicationError{Code: "101", Message: "Failed connection"}
	ResponseParseError        = ApplicationError{Code: "102", Message: "Failed parse response"}
	UnsupportedOperationError = ApplicationError{Code: "302", Message: "Unsupported Operation Type"}
)
