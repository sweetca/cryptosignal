package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/go-multierror"
	"github.com/sweetca/cryptosignal/common/model"
	"sync"
)

var ctx = context.Background()

type Service struct {
	client *redis.Client
}

func NewService(settings *Config) (*Service, error) {
	rClient := NewClient(settings)

	rService := Service{
		client: rClient,
	}
	if err := rService.Test(); err != nil {
		return nil, fmt.Errorf("error on redis service init: %v", err)
	}

	return &rService, nil
}

func (s *Service) Test() error {
	if err := s.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("fail redis ping: %v", err)
	}

	return nil
}

func (s *Service) StoreStatistic(statList []model.CryptoStatistic) error {
	var errList error
	var wg sync.WaitGroup

	for _, stat := range statList {
		key := stat.GetRedisKey()
		asBytes, err := stat.MarshalBinary()
		if err != nil {
			errList = multierror.Append(errList, err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.client.HSet(ctx, key, asBytes, 0).Err()
			if err != nil {
				errList = multierror.Append(errList, err)
			}
		}()
	}
	wg.Wait()

	return errList
}

func (s *Service) Publish(source string, message interface{}) error {
	err := s.client.Publish(ctx, source, message).Err()
	if err != nil {
		return fmt.Errorf("fail to send message to redis channel: %s : %v", source, err)
	}
	return nil
}
