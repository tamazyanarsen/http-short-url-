package data

import (
	"errors"
	"http-short-url/internal/config"
	file_handler "http-short-url/internal/fileHandler"
	"http-short-url/internal/logger"
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
	urlList sync.Map
}

func (rec *FileStore) InitMapStore() error {
	cons, consErr := file_handler.NewConsumer(*config.Config["f"])
	if consErr != nil {
		return consErr
	}

	data, dataErr := cons.ReadEvent()
	if dataErr != nil {
		return dataErr
	}

	rec.urlList.Store(data.ShortURL, data.OriginalURL)
	return nil
}

func (rec *FileStore) Write(key string, value string) error {
	rec.urlList.Store(key, value)
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
	if value, ok := rec.urlList.Load(key); ok {
		return value.(string), nil
	} else if cons, consErr := file_handler.NewConsumer(*config.Config["f"]); consErr == nil {
		for fileData, consErr := cons.ReadEvent(); consErr == nil; fileData, consErr = cons.ReadEvent() {
			if fileData.ShortURL == key {
				return fileData.OriginalURL, nil
			}
		}
	}

	return "", errors.New("false")
}
