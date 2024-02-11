package props

const KEY_SERVER_ENV = "KEY_SERVER_URL"

var KeyServerURL string

func loadKeyServerUrl() {
	KeyServerURL = readEnv(KEY_SERVER_ENV, "")
}
