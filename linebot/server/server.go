package server

import (
	"errors"
	"linebot/controller"
	"linebot/logger"
	"linebot/middle"
	"net/http"
	"sync"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

const (
	server_port  = ":1323"
	metrics_port = ":8081"
)

func requestLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Request(c.Request())
		return next(c)
	}
}

func StartServer() {
	serverGroup := new(sync.WaitGroup)
	// setup & run prometheus endpoints
	serverGroup.Add(1)
	go initMetricsServer(serverGroup)
	// setup & run app server
	serverGroup.Add(1)
	go initAppServer(serverGroup)
	logger.Info(logger.LBIF900001, server_port, metrics_port)
	serverGroup.Wait()
}

func initAppServer(sg *sync.WaitGroup) {

	controller := controller.NewLineEventController()

	appServer := echo.New()
	appServer.HideBanner = true
	appServer.HidePort = true
	appServer.Use(requestLog)
	appServer.Use(middle.VerifyLineSignature)
	appServer.Use(echoprometheus.NewMiddleware("linebot"))
	appServer.POST("/", controller.HandleLineAPIRequest)
	if err := appServer.Start(server_port); err != nil {
		logger.FatalWithStackTrace(err, logger.LBFT909999)
		sg.Done()
		panic(err)
	}
	sg.Done()
}

func initMetricsServer(sg *sync.WaitGroup) {
	metricsServer := echo.New()
	metricsServer.HideBanner = true
	metricsServer.HidePort = true
	metricsServer.GET("/metrics", echoprometheus.NewHandler())
	if err := metricsServer.Start(metrics_port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.FatalWithStackTrace(err, logger.LBFT909999)
		sg.Done()
		panic(err)
	}
	sg.Done()
}
