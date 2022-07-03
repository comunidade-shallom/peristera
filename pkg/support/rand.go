package support

import (
	"math/rand"
	"time"
)

func init() { //nolint:gochecknoinits
	rand.Seed(time.Now().UnixNano())
}

func RandBool() bool {
	return rand.Intn(2) == 1 //nolint:gomnd
}
