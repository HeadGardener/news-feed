package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBName     string
	Host       string
	SSLMode    string
	ServerPort string
}

func MustInit(path string) *Config {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("[FATAL] invalid path: %s", path)
		return nil
	}

	dbname := os.Getenv("dbname")
	if dbname == "" {
		log.Fatalf("[FATAL] db name is empty")
		return nil
	}

	host := os.Getenv("dbhost")
	if host == "" {
		log.Fatalf("[FATAL] db host is empty")
		return nil
	}

	sslmode := os.Getenv("sslmode")
	if sslmode == "" {
		log.Fatalf("[FATAL] sslmode is empty")
		return nil
	}

	srvport := os.Getenv("server_port")
	if srvport == "" {
		log.Fatalf("[FATAL] server_port is empty")
		return nil
	}

	return &Config{
		DBName:     dbname,
		Host:       host,
		SSLMode:    sslmode,
		ServerPort: srvport,
	}
}
