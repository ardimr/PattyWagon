package utils

import (
	"os"
	"strconv"
)

func GetEnvInt64(key string, defaultValue int64) int64 {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.ParseInt(key, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}
