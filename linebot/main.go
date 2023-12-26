package main

import (
	"linebot/dao"
	"linebot/logger"
	"linebot/server"
	"linebot/transfer"

	"github.com/labstack/echo/v4"
)

func requestLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Request(c.Request())
		return next(c)
	}
}

const (
	SERVER_PORT  = "1323"
	METRICS_PORT = "8081"
)

func main() {
	dao.InitDB()
	defer dao.Close()

	transfer.InitLineBot()

	server.StartServer()
}
