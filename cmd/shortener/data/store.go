package data

import (
	"sync"
)

type URLStore struct {
	urlList sync.Map
}

func (receiver *URLStore) Write(key string, value string) {
	println("URLStore WRITE key:value", key, value)
	receiver.urlList.Store(key, value)
}

func (receiver *URLStore) Read(key string) (string, bool) {
	println("URLStore READ key", key)
	if res, ok := receiver.urlList.Load(key); ok {
		return res.(string), true
	} else {
		return "", false
	}
}

type Store interface {
	Write(key string, value string)
	Read(key string) (string, bool)
}
