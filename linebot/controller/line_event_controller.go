package controller

import (
	"linebot/applicationerror"
	"linebot/dao/database"
	"linebot/dao/operation_history"
	"linebot/dao/user_info"
	"linebot/logger"
	"linebot/processor"
	"linebot/props"
	"linebot/security"
	"linebot/transfer/key_server"
	"linebot/transfer/line"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineEventController struct {
	opProcessor  *processor.OperationProcessor
	lineTransfer line.LineTransfer
}

func NewLineEventController() *LineEventController {
	// init Transfer
	lineBot, err := linebot.New(props.ChannelSecret, props.ChannelToken)

	if err != nil {
		logger.FatalWithStackTrace(err, applicationerror.SystemError, logger.LBFT040004)
		panic(err)
	}
	lTransfer := line.NewLineTransfer(lineBot)

	c := http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
		Timeout: 9 * time.Second,
	}
	ksTransfer := key_server.NewKeyServerTransferImpl(&c)

	// init DatabaseConnection
	database := &database.MySQLDatabaseConnection{}
	database.InitDB()
	ohDao := operation_history.OperationHistoryDaoImpl{Database: database}
	uiDao := user_info.UserInfoDaoImpl{Database: database}

	encryptor := security.EncryptorImpl{}
	lec := LineEventController{}
	lec.opProcessor = processor.NewOperationProcessor(
		ohDao,
		uiDao,
		lTransfer,
		ksTransfer,
		encryptor,
	)
	lec.lineTransfer = lTransfer
	return &lec
}

func (lec *LineEventController) HandleLineAPIRequest(c echo.Context) error {
	events, err := linebot.ParseRequest(props.ChannelSecret, c.Request())
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	lec.opProcessor.HandleEvents(events)

	return c.String(http.StatusOK, "")
}
