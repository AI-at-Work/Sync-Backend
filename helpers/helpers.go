package helpers

import (
	"context"
	"errors"
	"log"
)

func ContextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		log.Println("Request Canceled")
		return errors.New("Request Canceled")
	case context.DeadlineExceeded:
		log.Println("DeadLine Exceeded")
		return errors.New("DeadLine Exceeded")
	default:
		return nil
	}
}
