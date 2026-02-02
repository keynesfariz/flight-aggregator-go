package libs

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	godotenv.Load()
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
