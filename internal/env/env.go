package env

import (
	"os"
	"strconv"
)

// GetString reads environment variable, returns fallback If the key
// not exists, otherwise return value of the key.
func GetString(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}

// GetInt reads environment variable, then convert string value to integer
// otherwise, returns fallback value.
func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	intVal, err := strconv.Atoi(key)
	if err != nil {
		return fallback
	}

	return intVal
}

func GetFloat64(key string, fallback float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	floatVal, err := strconv.ParseFloat(key, 64)
	if err != nil {
		return fallback
	}

	return floatVal
}

func GetBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}
