package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"work-tracker/internal/model"
	"work-tracker/internal/secret"

	"github.com/redis/go-redis/v9"
)

type MonthRecapStore struct{ r redis.UniversalClient; sec *secret.Secret }

func NewMonthRecapStore(r redis.UniversalClient, sec *secret.Secret) *MonthRecapStore {
	return &MonthRecapStore{r: r, sec: sec}
}

func (s *MonthRecapStore) key(userID string, year int, month int) string {
	return fmt.Sprintf("monthrecap:%s:%04d-%02d", userID, year, month)
}

func (s *MonthRecapStore) Save(ctx context.Context, m *model.MonthRecap) error {
	b, _ := json.Marshal(m)
	enc, err := s.sec.Encrypt(b)
	if err != nil {
		return err
	}
	// configurable TTL; default 365 days
	return s.r.Set(ctx, s.key(m.UserID, m.Year, m.Month), enc, 365*24*time.Hour).Err()
}

func (s *MonthRecapStore) Get(ctx context.Context, userID string, year int, month int) (*model.MonthRecap, error) {
	val, err := s.r.Get(ctx, s.key(userID, year, month)).Result()
	if err != nil {
		return nil, err
	}
	pt, err := s.sec.DecryptString(val)
	if err != nil {
		return nil, err
	}
	var m model.MonthRecap
	if err := json.Unmarshal(pt, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
