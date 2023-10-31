package main

import (
	"errors"
	"fmt"
	"linebot/logger"
	"linebot/middle"
	"linebot/processer"
	"linebot/props"
	"net/http"
	"sync"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
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
	serverGroup := new(sync.WaitGroup)

	bot, err := linebot.New(props.ChannelSecret, props.ChannelToken)
	if err != nil {
		fmt.Println(err)
	}
	e := echo.New()
	e.Use(echoprometheus.NewMiddleware("linebot"))
	e.HideBanner = true
	e.HidePort = true
	e.Use(requestLog)
	e.Use(middle.VerifyLineSignature)
	e.POST("/", func(c echo.Context) error {
		events, err := bot.ParseRequest(c.Request())
		//_, err := bot.ParseRequest(c.Request())
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		processer.HandleEvents(bot, events)

		return c.String(http.StatusOK, "")
	})

	// setup & run prometheus endpoints
	serverGroup.Add(1)
	go func() {
		metrics := echo.New()
		metrics.HideBanner = true
		metrics.HidePort = true
		metrics.GET("/metrics", echoprometheus.NewHandler())
		if err := metrics.Start(fmt.Sprintf(":%s", METRICS_PORT)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			metrics.Logger.Fatal(err)
			serverGroup.Done()
		}
		serverGroup.Done()
	}()
	// setup & run app server
	serverGroup.Add(1)
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", SERVER_PORT)); err != nil {
			e.Logger.Fatal(e)
			serverGroup.Done()
		}
		serverGroup.Done()
	}()
	logger.Info(&logger.LBIF900001, SERVER_PORT, METRICS_PORT)
	serverGroup.Wait()
}
