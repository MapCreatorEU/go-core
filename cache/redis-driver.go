package cache

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisDriver struct {
	client *redis.Client
}

type RedisConfig struct {
	Addr string
	Password string
	DB int
}

func NewRedisDriver(Config RedisConfig)(*RedisDriver) {
	return &RedisDriver{
		client: redis.NewClient(&redis.Options{
			Addr: Config.Addr,
			Password: Config.Password,
			DB: Config.DB,
		}),
	}
}

func (Driver *RedisDriver) SetKeyWithDuration(Key string, Value string, expiration time.Duration) (error) {
	_, err := Driver.client.Set(Key, Value, expiration).Result()
	return err
}

func (Driver *RedisDriver) SetKey(Key string, Value string) (error) {
	return Driver.SetKeyWithDuration(Key, Value, 0)
}

func (Driver *RedisDriver) HasKey(Key string) (bool, error) {
	_, err := Driver.client.Get(Key).Result()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (Driver *RedisDriver) SAdd(Key string, Value string) (error) {
	_, err := Driver.client.SAdd(Key, Value).Result()
	return err
}

func (Driver *RedisDriver) SRem(Key string, Value string) (error) {
	_, err := Driver.client.SRem(Key, Value).Result()
	return err
}

func (Driver *RedisDriver) SHasKey(Key string, Value string) (bool, error) {
	return Driver.client.SIsMember(Key, Value).Result()
}

func (Driver *RedisDriver) Close() {
	Driver.client.Close()
}