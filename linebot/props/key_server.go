package props

import "net/url"

const KEY_SERVER_ENV = "KEY_SERVER_URL"

var KeyServerURL string

func loadKeyServerUrl() {
	loadUrl := readEnv(KEY_SERVER_ENV, "http://localhost:9999/")
	_, err := url.ParseRequestURI(loadUrl)
	if err != nil {
		panic("URL Environment is invalid: loadUrl")
	}
	KeyServerURL = loadUrl
}
