package key_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"linebot/applicationerror"
	"linebot/entity"
	"linebot/logger"
	"linebot/props"
	"net/http"

	"github.com/pkg/errors"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
type KeyServerTransferImpl struct {
	client HttpClient
}

func NewKeyServerTransferImpl(client HttpClient) *KeyServerTransferImpl {
	return &KeyServerTransferImpl{client}
}

func (kst KeyServerTransferImpl) Request(path string) (entity.KeyServerResponse, error) {
	logger.Info(logger.LBIF040001, path)
	req, _ := http.NewRequest("GET", props.KeyServerURL+path, bytes.NewReader([]byte("")))
	// res, err := kst.client.Get(props.KeyServerURL + path)
	res, err := kst.client.Do(req)
	logger.Debug(fmt.Sprintf("res%v", res))
	if err != nil {
		err = errors.Wrap(err, "Failed connect key server")
		logger.FatalWithStackTrace(err, logger.LBFT040001)
		return entity.KeyServerResponse{}, applicationerror.ConnectionError
	}

	bytesArray, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.Wrap(err, "Failed read response from key server")
		logger.FatalWithStackTrace(err, logger.LBFT040001)
		return entity.KeyServerResponse{}, applicationerror.ResponseParseError
	}
	res.Body.Close()
	logger.Info(logger.LBIF040002, string(bytesArray))
	var ret entity.KeyServerResponse

	err = json.Unmarshal(bytesArray, &ret)
	if err != nil {
		err = errors.Wrap(err, "Failed convert response from key server")
		logger.FatalWithStackTrace(err, logger.LBFT040003, string(bytesArray))
		return entity.KeyServerResponse{}, applicationerror.ResponseParseError
	}

	return ret, nil
}
func (kst KeyServerTransferImpl) OpenKey() (entity.KeyServerResponse, error) {
	return kst.Request("open")
}

func (kst KeyServerTransferImpl) CloseKey() (entity.KeyServerResponse, error) {
	return kst.Request("close")
}

func (kst KeyServerTransferImpl) CheckKey() (entity.KeyServerResponse, error) {
	return kst.Request("check")
}
