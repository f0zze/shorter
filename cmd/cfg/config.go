package cfg

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Host        string
	Response    string
	LogFilePath string
}

func GetConfig() ServerConfig {
	host := flag.String("a", "localhost:8080", "Server URL")
	destHost := flag.String("b", "http://localhost:8080", "Response server URL")
	flag.Parse()

	config := ServerConfig{
		*host,
		*destHost,
		"server.log",
	}

	if envRunAdd := os.Getenv("SERVER_ADDRESS"); envRunAdd != "" {
		config.Host = envRunAdd
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.Response = envBaseURL
	}

	return config
}
