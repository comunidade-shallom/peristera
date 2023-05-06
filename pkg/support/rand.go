package support

import (
	"crypto/rand"
)

func RandBool() bool {
	var randomInt int64

	randomBytes := make([]byte, 1)
	_, _ = rand.Read(randomBytes)
	randomInt = int64(randomBytes[0])

	return randomInt%2 == 0
}
