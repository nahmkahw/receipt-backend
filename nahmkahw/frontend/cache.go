package frontend

import (
	"encoding/json"
	"time"
	"fmt"

	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	redis_cache *redis.Client
	expires time.Duration
}

func NewRedisCacheFrontEnd(redis_cache *redis.Client, expires time.Duration) CacheFrontend {
	return &redisCache{redis_cache: redis_cache, expires: expires}
}

type CacheFrontend interface {
	SetOrderAll(key string, value *[]Order)
	GetOrderAll(key string) *[]Order

	SetYearSemesterAll(key string, value *[]YearSemester)
	GetYearSemesterAll(key string) *[]YearSemester

	SetLoginStatus(key string, value *Loginstatus)
	GetLoginStatus(key string) *Loginstatus
	DeleteLoginStatus(key string) 

}

func (cache *redisCache) DeleteLoginStatus(key string)  {
	val := cache.redis_cache.Del(key)
	fmt.Println(val)
}

func (cache *redisCache) SetLoginStatus(key string, value *Loginstatus) {

	json, err := json.Marshal(value)

	if err != nil {
		panic(err)
	}
	cache.redis_cache.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) GetLoginStatus(key string) *Loginstatus {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	login := Loginstatus{}
	err = json.Unmarshal([]byte(val), &login)
	if err != nil {
		panic(err)
	}

	return &login
}

func (cache *redisCache) SetOrderAll(key string, value *[]Order) {

	json, err := json.Marshal(value)

	if err != nil {
		panic(err)
	}
	cache.redis_cache.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) GetOrderAll(key string) *[]Order {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	order := []Order{}
	err = json.Unmarshal([]byte(val), &order)
	if err != nil {
		panic(err)
	}

	return &order
}

func (cache *redisCache) SetYearSemesterAll(key string, value *[]YearSemester) {

	json, err := json.Marshal(value)

	if err != nil {
		panic(err)
	}
	cache.redis_cache.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) GetYearSemesterAll(key string) *[]YearSemester {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	YearSemester := []YearSemester{}
	err = json.Unmarshal([]byte(val), &YearSemester)
	if err != nil {
		panic(err)
	}

	return &YearSemester
}
