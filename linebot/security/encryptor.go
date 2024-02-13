package security

import (
	"fmt"
	"linebot/props"

	"golang.org/x/crypto/scrypt"
)

type Encryptor interface {
	SaltHash(string) string
}

type EncryptorImpl struct{}

func (e EncryptorImpl) SaltHash(value string) string {
	ret, _ := scrypt.Key([]byte(value), []byte(props.Salt), 32768, 8, 1, 32)

	return fmt.Sprintf("%x", ret)
}
