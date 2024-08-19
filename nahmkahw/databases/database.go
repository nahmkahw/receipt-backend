package databases

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type connection struct{}

func NewDatabases() *connection {
	return &connection{}
}

func (c *connection) OracleInit() *sqlx.DB {

	db, err := oracleConnection()
	if err != nil {
		panic(err)
	}
	return db
}

func (c *connection) RedisInint() *redis.Client {
	return redisConnection()
}

func redisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		//Addr: viper.GetString("redis_cache.addressInLocal"),
		Addr: viper.GetString("redis_cache.address"),
		// Addr:     viper.GetString("redis_cache.addressInServer"),
		Password: viper.GetString("redis_cache.password"),
		DB:       viper.GetInt("redis_cache.db-num"),
	})
}

func oracleConnection() (*sqlx.DB, error) {

	dns := fmt.Sprintf("%v", viper.GetString("db.connection"))
	driver := viper.GetString("db.openDriver")

	return sqlx.Open(driver, dns)

}
