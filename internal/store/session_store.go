package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionStore struct{ r redis.UniversalClient }

func NewSessionStore(r redis.UniversalClient) *SessionStore { return &SessionStore{r: r} }

func (s *SessionStore) key(jti string) string { return fmt.Sprintf("session:%s", jti) }

func (s *SessionStore) Create(ctx context.Context, jti string, userID string, ttl time.Duration) error {
	return s.r.Set(ctx, s.key(jti), userID, ttl).Err()
}

func (s *SessionStore) Exists(ctx context.Context, jti string) (bool, error) {
	_, err := s.r.Get(ctx, s.key(jti)).Result()
	if err == redis.Nil {
		return false, nil
	}
	return err == nil, err
}

func (s *SessionStore) Delete(ctx context.Context, jti string) error {
	return s.r.Del(ctx, s.key(jti)).Err()
}
