package port

import "context"

type Customer interface {
	Get(context.Context, int) (int, error)
}
