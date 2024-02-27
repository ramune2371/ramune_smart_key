package applicationerror

import "fmt"

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

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("ApplicationError %s:%s", e.Code, e.Message)

}

func (e ApplicationError) ApplicationError() bool {
	return true
}

var (
	ConnectionError           = newApplicationError("101", "Failed connection")
	ResponseParseError        = newApplicationError("102", "Failed parse response")
	DBInsertError             = newApplicationError("201", "Failed DB Insert")
	DBUpdateError             = newApplicationError("202", "Failed DB Update")
	DBSelectError             = newApplicationError("203", "Failed DB Select")
	UnsupportedOperationError = newApplicationError("302", "Unsupported Operation Type")
	SignatureVerifyError      = newApplicationError("401", "LineWebhookEventの署名検証に失敗")
	ReplyError                = newApplicationError("402", "返信時にエラーが発生。")
	SystemError               = newApplicationError("999", "システムエラー。予期しないエラーが発生。")
)
