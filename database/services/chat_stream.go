package services

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

func (dataBase *Database) AckStream(ctx context.Context, val *redis.XMessage) {
	dataBase.Cache.XAck(ctx, os.Getenv("REDIS_STREAM"), os.Getenv("REDIS_STREAM_GROUP"), val.ID)
}

func (dataBase *Database) DelFromStream(ctx context.Context, val *redis.XMessage) {
	dataBase.Cache.XDel(ctx, os.Getenv("REDIS_STREAM"), val.ID)
}

func (dataBase *Database) ReadFromStream(ctx context.Context, count int64, timeout time.Duration, myKeyId string) *redis.XStreamSliceCmd {
	readGroup := dataBase.Cache.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    os.Getenv("REDIS_STREAM_GROUP"),
		Consumer: os.Getenv("REDIS_WORKER_NAME"),
		Streams:  []string{os.Getenv("REDIS_STREAM"), myKeyId},
		Count:    count,   // No Of Data To Retrieve
		Block:    timeout, //TimeOut
		NoAck:    false,
	})
	return readGroup

}
