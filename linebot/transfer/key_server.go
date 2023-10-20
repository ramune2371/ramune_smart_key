package transfer

import (
	"encoding/json"
	"fmt"
	"io"
	"linebot/entity"
	"linebot/logger"
	"net/http"
)

func Open() entity.KeyServerResponse {
  return request("open")
}

func Close() entity.KeyServerResponse {
  return request("close")
}

func Check() entity.KeyServerResponse {
  return request("check")
}


func request(path string) entity.KeyServerResponse {
  res,_ := http.Get("http://192.168.11.200:80/"+path)

  bytesArray,_ := io.ReadAll(res.Body)
  res.Body.Close()
  var ret entity.KeyServerResponse

  err := json.Unmarshal(bytesArray,&ret)
  if err != nil {
    logger.Debug(string(bytesArray))
    logger.Debug(err.Error())
  }

  logger.Debug(fmt.Sprintf("key server response = keyStatus:%s, operationStatus:%s",ret.KeyStatus,ret.OperationStatus))

  return ret
}
