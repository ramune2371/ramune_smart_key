package applicationerror

// TODO Goのエラー生成うまく理解してない。要勉強＆リファクタ

type ApplicationError struct {
	Code    string
	Message string
}

func newApplicationError(code string, message string) *ApplicationError {
	return &ApplicationError{
		Code:    code,
		Message: message,
	}
}

func (e *ApplicationError) Error() string { return e.Message }

func (e ApplicationError) ApplicationError() bool {
	return true
}

var (
	ConnectionError           = newApplicationError("101", "Failed connection")
	ResponseParseError        = newApplicationError("102", "Failed parse response")
	UnsupportedOperationError = newApplicationError("302", "Unsupported Operation Type")
)
