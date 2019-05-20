package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/gomodule/redigo/redis"
)

const (
	deviceStatusKeyPrefix = "ds"
)

var (
	_ store.DeviceStatuses = &Store{}
)

type Store struct {
	pool *redis.Pool
}

func NewStore(pool *redis.Pool) *Store {
	return &Store{
		pool: pool,
	}
}

func (s *Store) ResetDeviceStatus(ctx context.Context, deviceID string, ttl time.Duration) error {
	conn := s.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%s", deviceStatusKeyPrefix, deviceID)
	_, err := conn.Do("SETEX", key, ttl.Seconds(), "")
	return err
}

func (s *Store) GetDeviceStatus(ctx context.Context, deviceID string) (models.DeviceStatus, error) {
	conn := s.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%s", deviceStatusKeyPrefix, deviceID)
	_, err := redis.String(conn.Do("GET", key))
	switch err {
	case nil:
		return models.DeviceStatusOnline, nil
	case redis.ErrNil:
		return models.DeviceStatusOffline, nil
	default:
		return "", err
	}
}

func (s *Store) GetDeviceStatuses(ctx context.Context, deviceIDs []string) ([]models.DeviceStatus, error) {
	if len(deviceIDs) == 0 {
		return nil, nil
	}

	conn := s.pool.Get()
	defer conn.Close()

	var keys []interface{}
	for _, deviceID := range deviceIDs {
		keys = append(keys, fmt.Sprintf("%s:%s", deviceStatusKeyPrefix, deviceID))
	}

	vals, err := redis.Values(conn.Do("MGET", keys...))
	if err != nil {
		return nil, err
	}

	var deviceStatuses []models.DeviceStatus
	for _, val := range vals {
		if val == nil {
			deviceStatuses = append(deviceStatuses, models.DeviceStatusOffline)
		} else {
			deviceStatuses = append(deviceStatuses, models.DeviceStatusOnline)
		}
	}

	return deviceStatuses, nil
}
