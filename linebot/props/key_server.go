package props

import (
	"linebot/applicationerror"
	"linebot/logger"
	"net/url"
	"os"
)

const KEY_SERVER_ENV = "KEY_SERVER_URL"

var KeyServerURL string

// テストをしやすくるために、os.Exitを変数として定義
var OsExit = os.Exit

func loadKeyServerUrl() {
	loadUrl := readEnv(KEY_SERVER_ENV, "http://localhost:9999/")
	_, err := url.ParseRequestURI(loadUrl)
	if err != nil {
		logger.FatalWithStackTrace(err, applicationerror.SystemError, logger.LBFT900002, KEY_SERVER_ENV)
		OsExit(1)
	}
	KeyServerURL = loadUrl
}
