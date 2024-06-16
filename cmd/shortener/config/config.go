package config

import (
	"flag"
	"os"
)

var Config = make(map[string]*string)

func init() {
	Config["a"] = flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")
	Config["b"] = flag.String("b", "http://localhost:8080/", "базовый адрес результирующего сокращённого url")
	Config["f"] = flag.String("f", "/tmp/short-url-db.json", "полное имя файла, куда сохраняются данные в формате JSON (по умолчанию /tmp/short-url-db.json, пустое значение отключает функцию записи на диск)")
	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		Config["a"] = &serverAddr
	}
	if baseAddr := os.Getenv("BASE_URL"); baseAddr != "" {
		Config["b"] = &baseAddr
	}
	if fileLocation, exist := os.LookupEnv("FILE_STORAGE_PATH"); exist {
		Config["f"] = &fileLocation
	}
	println(*Config["f"])
}
