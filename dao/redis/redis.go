package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"), // Redis 密码
		DB:       viper.GetInt("redis.db"),          // Redis 数据库编号
		PoolSize: viper.GetInt("redis.pool_size"),   // 连接池大小
	})

	_, err = rdb.Ping().Result() // 通过 Ping 来测试连接
	fmt.Printf("connect redis success!\n")
	return
}
