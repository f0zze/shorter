package cfg

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Host            string
	Response        string
	LogFilePath     string
	FileStoragePath string
	DSN             string
}

func GetConfig() ServerConfig {
	host := flag.String("a", "localhost:8080", "Server URL")
	destHost := flag.String("b", "http://localhost:8080", "Response server URL")
	fileStoragePath := flag.String("f", "/tmp/short-url-db.json", "Url file storage path")
	dsn := flag.String("d", "", "Database DSN")
	flag.Parse()

	config := ServerConfig{
		*host,
		*destHost,
		"server.log",
		*fileStoragePath,
		*dsn,
	}

	if envRunAdd := os.Getenv("SERVER_ADDRESS"); envRunAdd != "" {
		config.Host = envRunAdd
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.Response = envBaseURL
	}

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		config.FileStoragePath = envFilePath
	}

	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		config.DSN = dsn
	}

	return config
}
