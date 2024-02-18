package security

import (
	"fmt"
	"linebot/logger"
	"linebot/props"
	"os"

	"golang.org/x/crypto/scrypt"
)

type Encryptor interface {
	SaltHash(string) string
}

type EncryptorImpl struct{}

func (e EncryptorImpl) SaltHash(value string) string {
	ret, err := scrypt.Key([]byte(value), []byte(props.Salt), 32768, 8, 1, 32)
	if err != nil {
		// scrypt.Keyの第3引数以降に固定値を入れているため、改修しない限り発生しない
		logger.FatalWithStackTrace(err, &logger.LBFT900001)
		os.Exit(1)
	}

	return fmt.Sprintf("%x", ret)
}
