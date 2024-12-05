package data

import (
	"http-short-url/internal/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(tesingObj *testing.T) {
	logger.InitLogger()
	var store Store = new(URLStore)
	store.Write("testKey", "testValue")
	mapValue, err := store.Read("testKey")
	tesingObj.Run("map-write", func(t *testing.T) {
		require.Nil(t, err)
		require.Equal(t, "testValue", mapValue)
	})
}
