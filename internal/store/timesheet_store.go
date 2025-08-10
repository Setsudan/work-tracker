package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"work-tracker/internal/model"

	"github.com/redis/go-redis/v9"
)

type TimeSheetStore struct{ r redis.UniversalClient }

func NewTimeSheetStore(r redis.UniversalClient) *TimeSheetStore { return &TimeSheetStore{r: r} }

func (s *TimeSheetStore) key(userID, weekStartISO string) string {
	return fmt.Sprintf("timesheet:%s:%s", userID, weekStartISO)
}

func (s *TimeSheetStore) Save(ctx context.Context, t *model.TimeSheet) error {
	b, _ := json.Marshal(t)
	// 90 days TTL
	return s.r.Set(ctx, s.key(t.UserID, t.WeekStartISO), b, 90*24*time.Hour).Err()
}

func (s *TimeSheetStore) Get(ctx context.Context, userID, weekStartISO string) (*model.TimeSheet, error) {
	val, err := s.r.Get(ctx, s.key(userID, weekStartISO)).Bytes()
	if err != nil {
		return nil, err
	}
	var t model.TimeSheet
	if err := json.Unmarshal(val, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
