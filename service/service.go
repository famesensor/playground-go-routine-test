package service

import (
	"context"
	"log"
	"time"

	"github.com/famesensor/playground-go-routine-test/port"
)

type Service interface {
	Get(context.Context, int) (int, error)
}

type service struct {
	Postgres port.Postgres
	Redis    port.Redis
}

func New(postgres port.Postgres, redis port.Redis) Service {
	return &service{
		Postgres: postgres,
		Redis:    redis,
	}
}

func (s *service) Get(ctx context.Context, id int) (int, error) {
	res, err := s.Redis.Get(ctx, id)
	if err != nil {
		return 0, err
	}

	if res == 0 {
		res, err = s.Postgres.Get(ctx, id)
		if err != nil {
			return 0, err
		}

		go func() {
			if err := s.Redis.Set(ctx, "id", res, 5*time.Minute); err != nil {
				log.Println("[SERVICE] Redis set failed: ", err)
			}
		}()
	}

	return res, nil
}
