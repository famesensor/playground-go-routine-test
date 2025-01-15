package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/famesensor/playground-go-routine-test/port"
)

type Service interface {
	Get(context.Context, int) (int, error)
	GetWithWait(context.Context, int) (int, error)
	GetWithWaitChannel(ctx context.Context, id int) (int, error)
}

type service struct {
	postgres    port.Postgres
	redis       port.Redis
	customer    port.Customer
	transaction port.Transaction
}

func New(postgres port.Postgres, redis port.Redis, customer port.Customer, transaction port.Transaction) Service {
	return &service{
		postgres:    postgres,
		redis:       redis,
		customer:    customer,
		transaction: transaction,
	}
}

func (s *service) Get(ctx context.Context, id int) (int, error) {
	res, err := s.redis.Get(ctx, id)
	if err != nil {
		return 0, err
	}

	if res == 0 {
		res, err = s.postgres.Get(ctx, id)
		if err != nil {
			return 0, err
		}

		go func() {
			if err := s.redis.Set(ctx, "id", res, 5*time.Minute); err != nil {
				log.Println("[SERVICE] Redis set failed: ", err)
			}
		}()
	}

	return res, nil
}

func (s *service) GetWithWait(ctx context.Context, id int) (int, error) {
	var (
		wg  sync.WaitGroup
		err error
		c   int
		t   int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		c, err = s.customer.Get(ctx, id)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t, err = s.transaction.Get(ctx, id)
	}()

	wg.Wait()
	if err != nil {
		return 0, err
	}

	return c + t, nil
}

func (s *service) GetWithWaitChannel(ctx context.Context, id int) (int, error) {
	var (
		wg   sync.WaitGroup
		c    = make(chan int, 1)
		errC = make(chan error, 1)
		t    = make(chan int, 1)
		errT = make(chan error, 1)
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		s.getCustomerChannel(ctx, id, c, errC)
	}()
	go func() {
		defer wg.Done()
		s.getTransactionChannel(ctx, id, t, errT)
	}()
	wg.Wait()

	if err := <-errC; err != nil {
		return 0, err
	}

	if err := <-errT; err != nil {
		return 0, err
	}

	return <-c + <-t, nil
}

func (s *service) getCustomerChannel(ctx context.Context, id int, ch chan int, errCh chan error) {
	res, err := s.customer.Get(ctx, id)
	ch <- res
	errCh <- err
}

func (s *service) getTransactionChannel(ctx context.Context, id int, ch chan int, errCh chan error) {
	res, err := s.transaction.Get(ctx, id)
	ch <- res
	errCh <- err
}
