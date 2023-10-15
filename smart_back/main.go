package main

import (
	"fmt"
	"net/http"
	"smart_back/dao"
	"smart_back/validator"
  "ramune/modules/logger"

	"github.com/labstack/echo/v4"
)

func logRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
    logger.Request(c.Request())
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.Use(logRequest)
	e.POST("/", func(c echo.Context) error {
		dao.GetAll()
		validator.RequestValidation(c)
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
	fmt.Println("Hello world")
}
