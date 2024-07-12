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
	userId := val.Values["userId"].(string)
	sessionId := val.Values["sessionId"].(string)
	sessionPrompt := val.Values["sessionPrompt"].(string)
	chats := val.Values["chats"].(string)
	chatsSummary := val.Values["chatsSummary"].(string)
	isNew := val.Values["isNew"].(string)
	balance := val.Values["balance"].(float64)

	dataContext, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute))
	errData := database.AddChat(dataContext, sessionId, sessionPrompt, chats, chatsSummary, isNew)
	if helpers.ContextError(dataContext) != nil || errData != nil {
		cancel()
		return errData
	}
	cancel()

	dataContext, cancel = context.WithDeadline(ctx, time.Now().Add(time.Minute))
	errData = database.UpdateBalance(dataContext, userId, balance)
	if helpers.ContextError(dataContext) != nil || errData != nil {
		cancel()
		return errData
	}

	fmt.Println("CHAT :  ", time.Now().UTC())
	cancel()
	return nil

}
