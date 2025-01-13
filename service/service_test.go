package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/famesensor/playground-go-routine-test/port/mock"
	"github.com/famesensor/playground-go-routine-test/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type serviceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	ctx          context.Context
	service      service.Service
	mockPostgres *mock.MockPostgres
	mockRedis    *mock.MockRedis
	wg           *sync.WaitGroup
}

func (s *serviceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ctx = context.Background()
	s.mockPostgres = mock.NewMockPostgres(s.ctrl)
	s.mockRedis = mock.NewMockRedis(s.ctrl)
	s.wg = &sync.WaitGroup{}
	s.service = service.New(s.mockPostgres, s.mockRedis)
}

func (s *serviceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (s *serviceTestSuite) TestGetSuccessAndSetFailed() {
	s.mockRedis.EXPECT().Get(s.ctx, 2).Return(0, nil)
	s.mockPostgres.EXPECT().Get(s.ctx, 2).Return(200, nil)

	// Use WaitGroup to synchronize the goroutine
	s.wg.Add(1)
	s.mockRedis.EXPECT().Set(s.ctx, "id", 200, 5*time.Minute).DoAndReturn(func(ctx context.Context, key string, value int, ttl time.Duration) error {
		defer s.wg.Done() // Mark the goroutine as done
		return errors.New("redis set error")
	})

	res, err := s.service.Get(s.ctx, 2)
	s.wg.Wait() // Wait for the goroutine to finish

	s.NoError(err)
	s.Equal(200, res)
}

func (s *serviceTestSuite) TestGetSuccess() {
	s.mockRedis.EXPECT().Get(s.ctx, 2).Return(0, nil)
	s.mockPostgres.EXPECT().Get(s.ctx, 2).Return(200, nil)

	// Use WaitGroup to synchronize the goroutine
	s.wg.Add(1)
	s.mockRedis.EXPECT().Set(s.ctx, "id", 200, 5*time.Minute).DoAndReturn(func(ctx context.Context, key string, value int, ttl time.Duration) error {
		defer s.wg.Done() // Mark the goroutine as done
		return nil
	})

	res, err := s.service.Get(s.ctx, 2)
	s.wg.Wait() // Wait for the goroutine to finish

	s.NoError(err)
	s.Equal(200, res)
}

// func TestService_Get(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockPostgres := mock.NewMockPostgres(ctrl)
// 	mockRedis := mock.NewMockRedis(ctrl)

// 	svc := service.New(mockPostgres, mockRedis)
// 	ctx := context.Background()

// 	var wg sync.WaitGroup

// 	t.Run("should return value from Redis", func(t *testing.T) {
// 		mockRedis.EXPECT().Get(ctx, 1).Return(100, nil)

// 		res, err := svc.Get(ctx, 1)

// 		assert.NoError(t, err)
// 		assert.Equal(t, 100, res)
// 	})

// 	t.Run("should return value from Postgres and update Redis", func(t *testing.T) {
// 		mockRedis.EXPECT().Get(ctx, 2).Return(0, nil)
// 		mockPostgres.EXPECT().Get(ctx, 2).Return(200, nil)

// 		// Use WaitGroup to synchronize the goroutine
// 		wg.Add(1)
// 		mockRedis.EXPECT().Set(ctx, "id", 200, 5*time.Minute).DoAndReturn(func(ctx context.Context, key string, value int, ttl time.Duration) error {
// 			defer wg.Done() // Mark the goroutine as done
// 			return nil
// 		})

// 		res, err := svc.Get(ctx, 2)
// 		wg.Wait() // Wait for the goroutine to finish

// 		assert.NoError(t, err)
// 		assert.Equal(t, 200, res)
// 	})

// 	t.Run("should handle Redis Get error", func(t *testing.T) {
// 		mockRedis.EXPECT().Get(ctx, 3).Return(0, errors.New("redis error"))

// 		res, err := svc.Get(ctx, 3)

// 		assert.Error(t, err)
// 		assert.Equal(t, 0, res)
// 	})

// 	t.Run("should handle Postgres Get error", func(t *testing.T) {
// 		mockRedis.EXPECT().Get(ctx, 4).Return(0, nil)
// 		mockPostgres.EXPECT().Get(ctx, 4).Return(0, errors.New("postgres error"))

// 		res, err := svc.Get(ctx, 4)

// 		assert.Error(t, err)
// 		assert.Equal(t, 0, res)
// 	})

// 	t.Run("should log error on Redis Set failure", func(t *testing.T) {
// 		mockRedis.EXPECT().Get(ctx, 5).Return(0, nil)
// 		mockPostgres.EXPECT().Get(ctx, 5).Return(500, nil)
// 		wg.Add(1)
// 		mockRedis.EXPECT().Set(ctx, "id", 500, 5*time.Minute).DoAndReturn(func(ctx context.Context, key string, value int, ttl time.Duration) error {
// 			defer wg.Done() // Mark the goroutine as done
// 			return errors.New("redis set error")
// 		})

// 		res, err := svc.Get(ctx, 5)
// 		wg.Wait() // Wait for the goroutine to finish

// 		assert.NoError(t, err)
// 		assert.Equal(t, 500, res)
// 	})
// }
