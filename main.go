package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"sync"
	"sync-worker/database/services"
	"sync-worker/worker"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Unable to load .env")
		return
	}

	checkBacklog, err := strconv.ParseBool(os.Getenv("REDIS_STREAM_CHECK_BACKLOG"))
	if err != nil {
		panic(err)
	}

	database := services.GetDataBase()
	log.Println("Database connected")

	if err := database.InitialiseRedisStreams(); err != nil {
		log.Println("Unable to initialise streams", err)
		return
	}
	log.Println("Streams Initialise successfully")

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go worker.SyncWithStream(context.TODO(), database, checkBacklog, &waitGroup)

	waitGroup.Wait()

}
