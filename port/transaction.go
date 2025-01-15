package port

import "context"

type Transaction interface {
	Get(context.Context, int) (int, error)
}
