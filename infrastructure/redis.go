package infrastructure

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// func NewRedisClient(addrs, password string, db int) *redis.ClusterClient {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	splitAddrs := strings.Split(addrs, ",")
// 	rdb := redis.NewClusterClient(&redis.ClusterOptions{
// 		Addrs:    splitAddrs,
// 		Password: password,
// 	})

//		if err := rdb.Ping(ctx).Err(); err != nil {
//			zap.L().Fatal("failed to ping redis", zap.Error(err))
//		}
//		return rdb
//	}
func NewRedisClient(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		zap.L().Fatal("failed to ping redis", zap.Error(err))
	}

	zap.L().Info("Connected to Redis successfully")
	return client
}
