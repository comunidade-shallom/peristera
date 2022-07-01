package support

import (
	"os"
	"strconv"
	"time"
)

// GetEnv from OS.
func GetEnv(key, def string) string {
	if val := os.Getenv(key); len(val) > 0 {
		return val
	}

	return def
}

// GetEnvInt from OS.
func GetEnvInt(key string, def int) (int, error) {
	if val := os.Getenv(key); len(val) > 0 {
		return strconv.Atoi(val)
	}

	return def, nil
}

// GetEnvDur from OS.
func GetEnvDur(key string, def time.Duration) (time.Duration, error) {
	if val := os.Getenv(key); len(val) > 0 {
		return time.ParseDuration(val)
	}

	return def, nil
}
