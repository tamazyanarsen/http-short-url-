package data

import (
	"errors"
	"http-short-url/cmd/shortener/config"
	file_handler "http-short-url/cmd/shortener/file-store"
	"http-short-url/cmd/shortener/logger"
	"strconv"
	"sync"
	"time"
)

type URLStore struct {
	urlList sync.Map
}

func (receiver *URLStore) Write(key string, value string) error {
	println("URLStore WRITE key:value", key, value)
	receiver.urlList.Store(key, value)
	return nil
}

func (receiver *URLStore) Read(key string) (string, error) {
	println("URLStore READ key", key)
	if res, ok := receiver.urlList.Load(key); ok {
		return res.(string), nil
	} else {
		return "", errors.New("false")
	}
}

type Store interface {
	Write(key string, value string) error
	Read(key string) (string, error)
}

type FileStore struct {
}

func (rec *FileStore) Write(key string, value string) error {
	prod, prodErr := file_handler.NewProducer(*config.Config["f"])
	if prodErr != nil {
		logger.Logger.Errorln(prodErr.Error())
		return prodErr
	}
	defer prod.Close()
	writeErr := prod.WriteEvent(&file_handler.FileData{
		UUID:        strconv.FormatInt(time.Now().Unix(), 10),
		ShortURL:    key,
		OriginalURL: value,
	})
	if writeErr != nil {
		logger.Logger.Errorln(writeErr.Error())
		return writeErr
	}
	return nil
}

func (rec *FileStore) Read(key string) (string, error) {
	if cons, err := file_handler.NewConsumer(*config.Config["f"]); err == nil {
		for fileData, err := cons.ReadEvent(); err == nil; fileData, err = cons.ReadEvent() {
			if fileData.ShortURL == key {
				return fileData.OriginalURL, nil
			}
		}
	}

	return "", errors.New("false")
}
