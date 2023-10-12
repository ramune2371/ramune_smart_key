package validator

import (
  "fmt"
  "github.com/labstack/echo/v4"
)

func RequestValidation(req echo.Context) error {
  var body []byte 
  req.Request().Body.Read(body)
  fmt.Println(body)
  return nil
}
