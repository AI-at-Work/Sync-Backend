package worker

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync"
	"sync-worker/database/services"
)

func SyncWithStream(ctx context.Context, database *services.Database, checkBacklog bool, group *sync.WaitGroup) {
	lastId := "0"
	defer group.Done()

	for {
		myKeyId := ">" //For Undelivered Ids So that Each Consumer Get Unique Id.
		if checkBacklog {
			myKeyId = lastId
		}

		readGroup := database.ReadFromStream(ctx, 10, 300, myKeyId)
		if readGroup.Err() != nil {
			if strings.Contains(readGroup.Err().Error(), "timeout") {
				fmt.Println("STREAM : TimeOUT")
				continue

			}
			if readGroup.Err() == redis.Nil {
				fmt.Println("STREAM : No Data Available")
				continue

			}
			panic(readGroup.Err())

		}

		data, err := readGroup.Result()
		if err != nil {
			panic(err)

		}

		if len(data[0].Messages) == 0 {
			checkBacklog = false
			fmt.Println("STREAM : Started Checking For New Messages ..!!")
			continue

		}

		var val redis.XMessage
		for _, val = range data[0].Messages {
			if val.Values["isNew"].(string) == "new" {
				err = SyncSessionDatabase(ctx, database, &val)
			}
			if err != nil {
				fmt.Println("STREAM SESSION : Error Please Check :", err)
				checkBacklog = true
				val.ID = "0"
				panic(err)
				//break
			}

			err = SyncChatDatabase(ctx, database, &val)
			if err != nil {
				fmt.Println("STREAM CHAT : Error Please Check :", err)
				checkBacklog = true
				val.ID = "0"
				panic(err)
				//break
			}

			fmt.Println("CHAT : Done Acknowledge")
			database.AckStream(ctx, &val)
			database.DelFromStream(ctx, &val)
		}
		lastId = val.ID
	}
}
