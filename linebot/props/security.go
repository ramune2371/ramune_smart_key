package props

const SALT_ENV = "SALT"

var Salt string

func loadSalt() {
	Salt = readEnv(SALT_ENV, "test")
}
