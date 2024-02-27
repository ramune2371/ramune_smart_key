package props

import "os"

func readEnv(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value != "" {
		return value
	} else {
		return defaultValue
	}
}

func LoadEnv() {
	loadLineProps()
	loadSalt()
	loadKeyServerUrl()
}
