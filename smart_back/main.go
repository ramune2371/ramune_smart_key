package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"smart_back/dao"
	"smart_back/validator"

	"github.com/labstack/echo/v4"
)

func logBody(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		buf, err := io.ReadAll(c.Request().Body)
		if err != nil {
			c.Logger().Error(err)
		}
		c.Logger().Info(string(buf))
		c.Request().Body = io.NopCloser(bytes.NewBuffer(buf))
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.Use(logBody)
	e.POST("/", func(c echo.Context) error {
		dao.GetAll()
		validator.RequestValidation(c)
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
	fmt.Println("Hello world")
}
