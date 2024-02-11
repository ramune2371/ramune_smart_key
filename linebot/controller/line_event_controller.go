package controller

import (
	"linebot/dao/database"
	"linebot/dao/operation_history"
	"linebot/dao/user_info"
	"linebot/processor"
	"linebot/security"
	"linebot/transfer/key_server"
	"linebot/transfer/line"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LineEventController struct {
	opProcessor  processor.OperationProcessor
	lineTransfer line.LineTransfer
}

func (lec *LineEventController) InitController() {
	// init Transfer
	lTransfer := new(line.LineTransferImpl)
	lTransfer.InitLineBot()
	ksTransfer := key_server.KeyServerTransferImpl{}

	// init DatabaseConnection
	database := &database.MySQLDatabaseConnection{}
	database.InitDB()
	ohDao := operation_history.OperationHistoryDaoImpl{Database: database}
	uiDao := user_info.UserInfoDaoImpl{Database: database}

	encryptor := security.EncryptorImpl{}

	lec.opProcessor = processor.OperationProcessor{OpHistoryDao: ohDao, UserInfoDao: uiDao, LineTransfer: lTransfer, KeyServerTransfer: ksTransfer, Encryptor: encryptor}
	lec.lineTransfer = lTransfer

}

func (lec *LineEventController) HandleLineAPIRequest(c echo.Context) error {
	events, err := lec.lineTransfer.ParseLineRequest(c.Request())
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	lec.opProcessor.HandleEvents(events)

	return c.String(http.StatusOK, "")
}
