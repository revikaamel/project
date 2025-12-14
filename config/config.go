package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	PGHost         string
	PGPort         string
	PGUser         string
	PGPassword     string
	PGDBName       string
	MongoURI       string
	MongoDB        string
	MongoCollection string
	JWTSecret      string
	JWTExpireHours int
	UploadDir      string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // ignore error if .env doesn't exist

	expire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))

	cfg := &Config{
		Port:           os.Getenv("PORT"),
		PGHost:         os.Getenv("PG_HOST"),
		PGPort:         os.Getenv("PG_PORT"),
		PGUser:         os.Getenv("PG_USER"),
		PGPassword:     os.Getenv("PG_PASSWORD"),
		PGDBName:       os.Getenv("PG_DBNAME"),
		MongoURI:       os.Getenv("MONGO_URI"),
		MongoDB:        os.Getenv("MONGO_DB"),
		MongoCollection: os.Getenv("MONGO_COLLECTION"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: expire,
		UploadDir:      os.Getenv("UPLOAD_DIR"),
	}

	log.Println("Configuration loaded")
	return cfg
}
