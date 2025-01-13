package port

import "context"

type Postgres interface {
	Get(context.Context, int) (int, error)
}
