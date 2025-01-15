package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/famesensor/playground-go-routine-test/port/mock"
	"github.com/famesensor/playground-go-routine-test/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type serviceTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	ctx             context.Context
	service         service.Service
	mockPostgres    *mock.MockPostgres
	mockRedis       *mock.MockRedis
	mockCustomer    *mock.MockCustomer
	mockTransaction *mock.MockTransaction
	wg              *sync.WaitGroup
}

func (s *serviceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ctx = context.Background()
	s.mockPostgres = mock.NewMockPostgres(s.ctrl)
	s.mockRedis = mock.NewMockRedis(s.ctrl)
	s.mockCustomer = mock.NewMockCustomer(s.ctrl)
	s.mockTransaction = mock.NewMockTransaction(s.ctrl)
	s.wg = new(sync.WaitGroup)
	s.service = service.New(s.mockPostgres, s.mockRedis, s.mockCustomer, s.mockTransaction)
}

func (s *serviceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

// func (s *serviceTestSuite) TestGetSuccessAndSetFailed() {
// 	s.mockRedis.EXPECT().Get(s.ctx, 2).Return(0, nil)
// 	s.mockPostgres.EXPECT().Get(s.ctx, 2).Return(200, nil)

// 	// Use WaitGroup to synchronize the goroutine
// 	s.wg.Add(1)
// 	s.mockRedis.EXPECT().Set(s.ctx, "id", 200, 5*time.Minute).DoAndReturn(func(ctx context.Context, key string, value int, ttl time.Duration) error {
// 		defer s.wg.Done() // Mark the goroutine as done
// 		return errors.New("redis set error")
// 	})

// 	res, err := s.service.Get(s.ctx, 2)
// 	s.wg.Wait() // Wait for the goroutine to finish

// 	s.NoError(err)
// 	s.Equal(200, res)
// }

// func (s *serviceTestSuite) TestGetSuccess() {
// 	s.mockRedis.EXPECT().Get(s.ctx, 2).Return(0, nil)
// 	s.mockPostgres.EXPECT().Get(s.ctx, 2).Return(200, nil)

// 	// Use WaitGroup to synchronize the goroutine
// 	s.wg.Add(1)
// 	s.mockRedis.EXPECT().Set(s.ctx, "id", 200, 5*time.Minute).DoAndReturn(func(ctx context.Context, key string, value int, ttl time.Duration) error {
// 		defer s.wg.Done() // Mark the goroutine as done
// 		return nil
// 	})

// 	res, err := s.service.Get(s.ctx, 2)
// 	s.wg.Wait() // Wait for the goroutine to finish

// 	s.NoError(err)
// 	s.Equal(200, res)
// }

// func (s *serviceTestSuite) TestGetWithWaitSuccess() {
// 	s.mockCustomer.EXPECT().Get(s.ctx, 2).Return(1, nil)
// 	s.mockTransaction.EXPECT().Get(s.ctx, 2).Return(1, nil)

// 	res, err := s.service.GetWithWait(s.ctx, 2)

// 	s.NoError(err)
// 	s.Equal(2, res)
// }

// func (s *serviceTestSuite) TestGetWithWaitFailed() {
// 	s.mockCustomer.EXPECT().Get(s.ctx, 2).Return(1, errors.New("redis set error"))
// 	s.mockTransaction.EXPECT().Get(s.ctx, 2).Return(1, nil)

// 	_, err := s.service.GetWithWait(s.ctx, 2)

// 	s.Error(err)
// }

func (s *serviceTestSuite) TestGetWithWaitChannelSuccess() {
	s.mockCustomer.EXPECT().Get(s.ctx, 2).Return(1, nil)
	s.mockTransaction.EXPECT().Get(s.ctx, 2).Return(2, nil)

	res, err := s.service.GetWithWaitChannel(s.ctx, 2)

	s.NoError(err)
	s.Equal(3, res)
}

func (s *serviceTestSuite) TestGetWithWaitChannelCustomerFailed() {
	s.mockCustomer.EXPECT().Get(s.ctx, 2).Return(1, errors.New("error"))
	s.mockTransaction.EXPECT().Get(s.ctx, 2).Return(2, nil)

	_, err := s.service.GetWithWaitChannel(s.ctx, 2)

	s.Error(err)
}

func (s *serviceTestSuite) TestGetWithWaitChannelTransactionFailed() {
	s.mockCustomer.EXPECT().Get(s.ctx, 2).Return(1, nil)
	s.mockTransaction.EXPECT().Get(s.ctx, 2).Return(0, errors.New("error"))

	_, err := s.service.GetWithWaitChannel(s.ctx, 2)

	s.Error(err)
}
