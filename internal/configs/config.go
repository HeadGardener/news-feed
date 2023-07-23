package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
}

type DBConfig struct {
	DBName  string
	Host    string
	SSLMode string
}

type ServerConfig struct {
	ServerPort string
}

func MustInit(path string) *Config {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("[FATAL] invalid path: %s", path)
		return nil
	}

	dbname := os.Getenv("DBNAME")
	if dbname == "" {
		log.Fatalf("[FATAL] db name is empty")
		return nil
	}

	host := os.Getenv("DBHOST")
	if host == "" {
		log.Fatalf("[FATAL] db host is empty")
		return nil
	}

	sslmode := os.Getenv("SSLMODE")
	if sslmode == "" {
		log.Fatalf("[FATAL] sslmode is empty")
		return nil
	}

	srvport := os.Getenv("SERVER_PORT")
	if srvport == "" {
		log.Fatalf("[FATAL] server port is empty")
		return nil
	}

	return &Config{
		DBConfig: DBConfig{
			DBName:  dbname,
			Host:    host,
			SSLMode: sslmode,
		},
		ServerConfig: ServerConfig{
			ServerPort: srvport,
		},
	}
}
