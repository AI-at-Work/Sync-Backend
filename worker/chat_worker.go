package worker

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync-worker/database/services"
	"sync-worker/helpers"
	"time"
)

func SyncChatDatabase(ctx context.Context, database *services.Database, val *redis.XMessage) error {
	sessionId := val.Values["sessionId"].(string)
	sessionPrompt := val.Values["sessionPrompt"].(string)
	chats := val.Values["chats"].(string)
	isNew := val.Values["isNew"].(string)

	dataContext, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute))

	errData := database.AddChat(dataContext, sessionId, sessionPrompt, chats, isNew)

	if helpers.ContextError(dataContext) != nil || errData != nil {
		cancel()
		return errData
	}

	fmt.Println("CHAT :  ", time.Now().UTC())
	cancel()
	return nil

}
