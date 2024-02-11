package security

import (
	"fmt"
	"linebot/logger"
	"linebot/props"

	"golang.org/x/crypto/scrypt"
)

type Encryptor interface {
	SaltHash(string) string
}

type EncryptorImpl struct{}

func (e EncryptorImpl) SaltHash(value string) string {
	ret, err := scrypt.Key([]byte(value), []byte(props.Salt), 32768, 8, 1, 32)
	if err != nil {
		logger.FatalWithStackTrace(err, &logger.LBFT900001)
	}
	return fmt.Sprintf("%x", ret)
}
