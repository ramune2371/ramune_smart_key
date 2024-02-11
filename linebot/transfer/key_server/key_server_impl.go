package key_server

import (
	"encoding/json"
	"io"
	"linebot/applicationerror"
	"linebot/entity"
	"linebot/logger"
	"linebot/props"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type KeyServerTransferImpl struct{}

func request(path string) (entity.KeyServerResponse, error) {
	logger.Info(&logger.LBIF040001, path)
	c := http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
		Timeout: 9 * time.Second,
	}

	res, err := c.Get(props.KeyServerURL + path)
	if err != nil {
		err = errors.Wrap(err, "Failed connect key server")
		logger.FatalWithStackTrace(err, &logger.LBFT040001)
		return entity.KeyServerResponse{}, &applicationerror.ConnectionError
	}

	bytesArray, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.Wrap(err, "Failed read response from key server")
		logger.FatalWithStackTrace(err, &logger.LBFT040001)
		return entity.KeyServerResponse{}, &applicationerror.ResponseParseError
	}
	res.Body.Close()
	logger.Info(&logger.LBIF040002, string(bytesArray))
	var ret entity.KeyServerResponse

	err = json.Unmarshal(bytesArray, &ret)
	if err != nil {
		err = errors.Wrap(err, "Failed convert response from key server")
		logger.FatalWithStackTrace(err, &logger.LBFT040003, string(bytesArray))
		return entity.KeyServerResponse{}, &applicationerror.ResponseParseError
	}

	return ret, nil
}
func (kt KeyServerTransferImpl) OpenKey() (entity.KeyServerResponse, error) {
	return request("open")
}

func (kt KeyServerTransferImpl) CloseKey() (entity.KeyServerResponse, error) {
	return request("close")
}

func (kt KeyServerTransferImpl) CheckKey() (entity.KeyServerResponse, error) {
	return request("check")
}
