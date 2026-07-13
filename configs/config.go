package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configs struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTexpiration int64
}

func Load() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system environment variables")
	}

	return &Configs{
		Port: os.Getenv("APP_PORT"),

		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSLMODE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTexpiration: getEnvInt("JWT_EXPIRATION", 3600*24*7),
	}
}
func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 14)
		if err != nil {
			return fallback
		}
		return i

	}
	return fallback

}
