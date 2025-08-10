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

type UserStore struct{ r redis.UniversalClient }

func NewUserStore(r redis.UniversalClient) *UserStore { return &UserStore{r: r} }

func (s *UserStore) userKey(id string) string     { return fmt.Sprintf("user:%s", id) }
func (s *UserStore) emailKey(email string) string { return fmt.Sprintf("user:email:%s", email) }
func (s *UserStore) usersSetKey() string          { return "users:all" }

func (s *UserStore) Create(ctx context.Context, u *model.User) error {
	id := ulid.Make().String()
	now := time.Now().UTC()
	u.ID = id
	u.CreatedAt = now
	u.UpdatedAt = now
	b, _ := json.Marshal(u)

	pipe := s.r.TxPipeline()
	pipe.Set(ctx, s.userKey(id), b, 0)
	pipe.Set(ctx, s.emailKey(u.Email), id, 0)
	pipe.SAdd(ctx, s.usersSetKey(), id)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	id, err := s.r.Get(ctx, s.emailKey(email)).Result()
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*model.User, error) {
	val, err := s.r.Get(ctx, s.userKey(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var u model.User
	if err := json.Unmarshal(val, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) Update(ctx context.Context, u *model.User) error {
	u.UpdatedAt = time.Now().UTC()
	b, _ := json.Marshal(u)
	return s.r.Set(ctx, s.userKey(u.ID), b, 0).Err()
}

func (s *UserStore) ListAllIDs(ctx context.Context) ([]string, error) {
	return s.r.SMembers(ctx, s.usersSetKey()).Result()
}
