package util

import (
	"encoding/json"
	"github.com/allegro/bigcache"
	"time"
)

var cache *bigcache.BigCache

func init() {
	var err error
	cache, err = bigcache.NewBigCache(bigcache.DefaultConfig(24 * time.Hour))
	if err != nil {
		panic("initialize cache middleware error, " + err.Error())
	}
}

func SetCache(key string, value interface{}) error {
	js, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cache.Set(key, js)
}

func GetCache(key string, obj interface{}) error {
	bt, err := cache.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, obj)
}

func ClearCache(key string) error {
	return cache.Delete(key)
}

func ClearAllCache() error {
	return cache.Reset()
}
