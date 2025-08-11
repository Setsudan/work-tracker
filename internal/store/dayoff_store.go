package store

import (
	"context"
	"encoding/json"
	"fmt"

	"work-tracker/internal/model"
	"work-tracker/internal/secret"

	"github.com/redis/go-redis/v9"
)

type DaysOffStore struct{ r redis.UniversalClient; sec *secret.Secret }

func NewDaysOffStore(r redis.UniversalClient, sec *secret.Secret) *DaysOffStore { return &DaysOffStore{r: r, sec: sec} }

func (s *DaysOffStore) key(userID string) string { return fmt.Sprintf("daysoff:%s", userID) }
func (s *DaysOffStore) itemKey(userID, dateISO string) string {
	return fmt.Sprintf("dayoff:%s:%s", userID, dateISO)
}

func (s *DaysOffStore) Add(ctx context.Context, userID string, d model.DayOff) error {
	b, _ := json.Marshal(d)
	enc, err := s.sec.Encrypt(b)
	if err != nil {
		return err
	}
	pipe := s.r.TxPipeline()
	pipe.SAdd(ctx, s.key(userID), d.DateISO)
	pipe.Set(ctx, s.itemKey(userID, d.DateISO), enc, 0)
	_, err = pipe.Exec(ctx)
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
		val, err := s.r.Get(ctx, s.itemKey(userID, date)).Result()
		if err != nil {
			continue
		}
		pt, err := s.sec.DecryptString(val)
		if err != nil {
			continue
		}
		var d model.DayOff
		if err := json.Unmarshal(pt, &d); err == nil {
			res = append(res, d)
		}
	}
	return res, nil
}
