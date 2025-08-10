package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"work-tracker/internal/model"

	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
)

type TimeLogStore struct{ r redis.UniversalClient }

func NewTimeLogStore(r redis.UniversalClient) *TimeLogStore { return &TimeLogStore{r: r} }

func (s *TimeLogStore) key(userID string) string { return fmt.Sprintf("timelogs:%s", userID) }

func (s *TimeLogStore) Add(ctx context.Context, userID string, logType model.TimeLogType, ts time.Time) (*model.TimeLog, error) {
	l := model.TimeLog{ID: ulid.Make().String(), UserID: userID, Type: logType, Timestamp: ts.UTC()}
	b, _ := json.Marshal(l)
	score := float64(l.Timestamp.UnixMilli())
	pipe := s.r.TxPipeline()
	pipe.ZAdd(ctx, s.key(userID), redis.Z{Score: score, Member: string(b)})
	cutoff := float64(time.Now().Add(-14 * 24 * time.Hour).UnixMilli())
	pipe.ZRemRangeByScore(ctx, s.key(userID), "-inf", fmt.Sprintf("%f", cutoff))
	_, err := pipe.Exec(ctx)
	return &l, err
}

func (s *TimeLogStore) GetLast(ctx context.Context, userID string) (*model.TimeLog, error) {
	vals, err := s.r.ZRevRange(ctx, s.key(userID), 0, 0).Result()
	if err != nil || len(vals) == 0 {
		return nil, err
	}
	var l model.TimeLog
	if err := json.Unmarshal([]byte(vals[0]), &l); err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *TimeLogStore) Range(ctx context.Context, userID string, from, to time.Time) ([]model.TimeLog, error) {
	min := fmt.Sprintf("%d", from.UTC().UnixMilli())
	max := fmt.Sprintf("%d", to.UTC().UnixMilli())
	vals, err := s.r.ZRangeByScore(ctx, s.key(userID), &redis.ZRangeBy{Min: min, Max: max}).Result()
	if err != nil {
		return nil, err
	}
	res := make([]model.TimeLog, 0, len(vals))
	for _, v := range vals {
		var l model.TimeLog
		if err := json.Unmarshal([]byte(v), &l); err == nil {
			res = append(res, l)
		}
	}
	return res, nil
}
