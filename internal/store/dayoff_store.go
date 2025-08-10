package store

import (
	"context"
	"encoding/json"
	"fmt"

	"work-tracker/internal/model"

	"github.com/redis/go-redis/v9"
)

type DaysOffStore struct{ r redis.UniversalClient }

func NewDaysOffStore(r redis.UniversalClient) *DaysOffStore { return &DaysOffStore{r: r} }

func (s *DaysOffStore) key(userID string) string { return fmt.Sprintf("daysoff:%s", userID) }
func (s *DaysOffStore) itemKey(userID, dateISO string) string {
	return fmt.Sprintf("dayoff:%s:%s", userID, dateISO)
}

func (s *DaysOffStore) Add(ctx context.Context, userID string, d model.DayOff) error {
	b, _ := json.Marshal(d)
	pipe := s.r.TxPipeline()
	pipe.SAdd(ctx, s.key(userID), d.DateISO)
	pipe.Set(ctx, s.itemKey(userID, d.DateISO), b, 0)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *DaysOffStore) Remove(ctx context.Context, userID, dateISO string) error {
	pipe := s.r.TxPipeline()
	pipe.SRem(ctx, s.key(userID), dateISO)
	pipe.Del(ctx, s.itemKey(userID, dateISO))
	_, err := pipe.Exec(ctx)
	return err
}

func (s *DaysOffStore) List(ctx context.Context, userID string) ([]model.DayOff, error) {
	dates, err := s.r.SMembers(ctx, s.key(userID)).Result()
	if err != nil {
		return nil, err
	}
	res := make([]model.DayOff, 0, len(dates))
	for _, date := range dates {
		val, err := s.r.Get(ctx, s.itemKey(userID, date)).Bytes()
		if err != nil {
			continue
		}
		var d model.DayOff
		if err := json.Unmarshal(val, &d); err == nil {
			res = append(res, d)
		}
	}
	return res, nil
}
