package worker

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync-worker/database/services"
	"sync-worker/helpers"
	"time"
)

func SyncSessionDatabase(ctx context.Context, database *services.Database, val *redis.XMessage) error {
	userId := val.Values["userId"].(string)
	sessionId := val.Values["sessionId"].(string)
	modelId := val.Values["modelId"].(string)
	sessionName := val.Values["sessionName"].(string)

	fmt.Println("SESSION : ", userId)
	fmt.Println("SESSION : ", sessionId)
	fmt.Println("SESSION : ", modelId)

	dataContext, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute))

	errData := database.AddSession(dataContext, userId, sessionId, modelId, sessionName)

	if errData != nil && helpers.ContextError(dataContext) != nil && !strings.Contains(errData.Error(), "duplicate key") {
		cancel()
		return errData
	}

	cancel()
	return nil

}
