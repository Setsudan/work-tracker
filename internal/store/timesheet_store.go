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

type TimeSheetStore struct{ r redis.UniversalClient; sec *secret.Secret }

func NewTimeSheetStore(r redis.UniversalClient, sec *secret.Secret) *TimeSheetStore {
	return &TimeSheetStore{r: r, sec: sec}
}

func (s *TimeSheetStore) key(userID, weekStartISO string) string {
	return fmt.Sprintf("timesheet:%s:%s", userID, weekStartISO)
}

func (s *TimeSheetStore) Save(ctx context.Context, t *model.TimeSheet) error {
	b, _ := json.Marshal(t)
	enc, err := s.sec.Encrypt(b)
	if err != nil {
		return err
	}
	// 90 days TTL
	return s.r.Set(ctx, s.key(t.UserID, t.WeekStartISO), enc, 90*24*time.Hour).Err()
}

func (s *TimeSheetStore) Get(ctx context.Context, userID, weekStartISO string) (*model.TimeSheet, error) {
	val, err := s.r.Get(ctx, s.key(userID, weekStartISO)).Result()
	if err != nil {
		return nil, err
	}
	pt, err := s.sec.DecryptString(val)
	if err != nil {
		return nil, err
	}
	var t model.TimeSheet
	if err := json.Unmarshal(pt, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
