package controller

import (
	"linebot/dao/database"
	"linebot/dao/operation_history"
	"linebot/dao/user_info"
	"linebot/processor"
	"linebot/transfer"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LineEventController struct {
	opProcessor processor.OperationProcessor
}

func (lec *LineEventController) InitController() {
	database := &database.MySQLDatabaseConnection{}
	database.InitDB()
	ohDao := operation_history.OperationHistoryDaoImpl{Database: database}
	uiDao := user_info.UserInfoDaoImpl{Database: database}
	lec.opProcessor = processor.OperationProcessor{OpHistoryDao: ohDao, UserInfoDao: uiDao}
}

func (lec *LineEventController) HandleLineAPIRequest(c echo.Context) error {
	events, err := transfer.ParseLineRequest(c.Request())
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	lec.opProcessor.HandleEvents(events)

	return c.String(http.StatusOK, "")
}
