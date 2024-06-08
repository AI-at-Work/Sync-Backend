package services

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
	"sync-worker/database/initialize"
)

type Database struct {
	Db    *sqlx.DB
	Cache *redis.Client
}

func GetDataBase() *Database {
	return &Database{
		Db:    initialize.InitPostgres(),
		Cache: initialize.InitRedis(),
	}
}

func (dataBase *Database) InitialiseRedisStreams() error {
	if err := dataBase.InitStream(context.Background()); err != nil {
		return err
	}
	return nil
}

func (dataBase *Database) InitStream(ctx context.Context) error {
	err := dataBase.Cache.XGroupCreateMkStream(ctx, os.Getenv("REDIS_STREAM"), os.Getenv("REDIS_STREAM_GROUP"), "$").Err()
	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			fmt.Println("Session Stream Already Exist")
			return nil
		}
		return err
	}
	return nil
}
