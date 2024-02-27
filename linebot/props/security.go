package props

import "strconv"

const (
	SALT_ENV        = "SALT"
	SIGNATURE_CHECK = "SIGNATURE_CHECK"
)

var (
	Salt             string
	IsSignatureCheck bool
)

func loadSalt() {
	Salt = readEnv(SALT_ENV, "test")
	signatureCheckStr := readEnv(SIGNATURE_CHECK, "true")
	isCheck, err := strconv.ParseBool(signatureCheckStr)
	if err != nil {
		IsSignatureCheck = true
	} else {
		IsSignatureCheck = isCheck
	}
}
