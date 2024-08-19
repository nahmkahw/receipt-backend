package backend

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	redis_cache *redis.Client
	expires time.Duration
}

func NewRedisCache(redis_cache *redis.Client, expires time.Duration) Cache{
	return &redisCache{redis_cache: redis_cache, expires: expires}
}

type Cache interface {
	SetStudent(key string, value *Student)
	GetStudent(key string) *Student

	SetOrder(key string, value *Order)
	GetOrder(key string) *Order

	SetOrderAll(key string, value *[]Order)
	GetOrderAll(key string) *[]Order

	SetPaymentAll(key string, value *[]Payment)
	GetPaymentAll(key string) *[]Payment

	SetReceiptAll(key string, value *[]Receipt)
	GetReceiptAll(key string) *[]Receipt
}

func (cache *redisCache) SetStudent(key string, value *Student) {

	json, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
	}
	cache.redis_cache.Set(key, json, cache.expires)
}

func (cache *redisCache) GetStudent(key string) *Student {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	student := Student{}
	err = json.Unmarshal([]byte(val), &student)
	if err != nil {
		fmt.Println(err)
	}

	return &student
}

func (cache *redisCache) SetOrder(key string, value *Order) {

	json, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
	}
	cache.redis_cache.Set(key, json, cache.expires)
}

func (cache *redisCache) GetOrder(key string) *Order {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	order := Order{}
	err = json.Unmarshal([]byte(val), &order)
	if err != nil {
		fmt.Println(err)
	}

	return &order
}

func (cache *redisCache) SetOrderAll(key string, value *[]Order) {

	json, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
	}

	cache.redis_cache.Set(key, json, cache.expires)
}

func (cache *redisCache) GetOrderAll(key string) *[]Order {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	order := []Order{}
	err = json.Unmarshal([]byte(val), &order)
	if err != nil {
		fmt.Println(err)
	}

	return &order
}

func (cache *redisCache) SetPaymentAll(key string, value *[]Payment) {

	json, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
	}
	cache.redis_cache.Set(key, json, cache.expires)
}

func (cache *redisCache) GetPaymentAll(key string) *[]Payment {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	payments := []Payment{}
	err = json.Unmarshal([]byte(val), &payments)
	if err != nil {
		fmt.Println(err)
	}

	return &payments
}

func (cache *redisCache) SetReceiptAll(key string, value *[]Receipt) {

	json, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
	}
	cache.redis_cache.Set(key, json, cache.expires)
}

func (cache *redisCache) GetReceiptAll(key string) *[]Receipt {

	val, err := cache.redis_cache.Get(key).Result()
	if err != nil {
		return nil
	}

	receipts := []Receipt{}
	err = json.Unmarshal([]byte(val), &receipts)
	if err != nil {
		fmt.Println(err)
	}

	return &receipts
}
