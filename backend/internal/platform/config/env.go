package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	_ = godotenv.Load()
}

func StringFromEnv(name string, current string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return current
}

func IntFromEnv(name string, current int) (int, error) {
	value := os.Getenv(name)
	if value == "" {
		return current, nil
	}
	return strconv.Atoi(value)
}

func DurationFromEnv(name string, current time.Duration) (time.Duration, error) {
	value := os.Getenv(name)
	if value == "" {
		return current, nil
	}
	return time.ParseDuration(value)
}
