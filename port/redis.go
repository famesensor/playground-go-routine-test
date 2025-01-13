package port

import (
	"context"
	"time"
)

type Redis interface {
	Set(context.Context, string, int, time.Duration) error
	Get(context.Context, int) (int, error)
}
