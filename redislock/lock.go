package redislock

import (
	"context"
	"errors"
	_ "embed"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	rdb  *redis.Client
	once sync.Once
)

// Init 必须在 main 中调用一次
func Init(addr string, db int, poolSize int) {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			DB:       db,
			PoolSize: poolSize,
		})
	})
}

// GetRDB 获取全局唯一的 Redis 客户端
func GetRDB() *redis.Client {
	if rdb == nil {
		panic("redislock not initialized")
	}
	return rdb
}

type RedisLock struct {
	rdb    *redis.Client
	key    string
	value  string // 唯一标识，用于解锁
	expire time.Duration
}

// NewRedisLock 创建一个分布式锁
func NewRedisLock(rdb *redis.Client, key string, expire time.Duration) *RedisLock {
	return &RedisLock{
		rdb:    rdb,
		key:    key,
		expire: expire,
	}
}

// Lock 尝试获取锁
// 返回 true 表示成功获取，false 表示已被占用
// 调用方需保存返回的 redislock 对象，并在业务结束后调用 Unlock()
func (rl *RedisLock) Lock(ctx context.Context) (bool, error) {
	// 生成唯一 value（必须全局唯一）
	rl.value = uuid.New().String()

	// SET key value NX PX milliseconds
	ok, err := rl.rdb.SetNX(ctx, rl.key, rl.value, rl.expire).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// unlockScript 是原子删除脚本：只有 value 匹配才删除 lua脚本，用embed编译时导入
//go:embed unlock.lua
var unlockScript string

// Unlock 释放锁
func (rl *RedisLock) Unlock(ctx context.Context) error {
	if rl.value == "" {
		return errors.New("redislock not acquired")
	}
	script := redis.NewScript(unlockScript)
	_, err := script.Run(ctx, rl.rdb, []string{rl.key}, rl.value).Result()
	if err != nil {
		return err
	}
	return nil
}
