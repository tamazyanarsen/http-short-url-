package config

import (
	"flag"
	"os"
)

var Config = make(map[string]*string)

func init() {
	Config["a"] = flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")
	Config["b"] = flag.String("b", "http://localhost:8080/", "базовый адрес результирующего сокращённого url")
	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		Config["a"] = &serverAddr
	}
	if baseAddr := os.Getenv("BASE_URL"); baseAddr != "" {
		Config["b"] = &baseAddr
	}
}
