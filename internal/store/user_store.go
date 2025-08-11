package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"work-tracker/internal/model"
	"work-tracker/internal/secret"

	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
)

type UserStore struct{ r redis.UniversalClient; sec *secret.Secret }

func NewUserStore(r redis.UniversalClient, sec *secret.Secret) *UserStore { return &UserStore{r: r, sec: sec} }

func (s *UserStore) userKey(id string) string { return fmt.Sprintf("user:%s", id) }
func (s *UserStore) emailKey(email string) string {
	// HMAC the email so it is not stored in plaintext keys
	return fmt.Sprintf("user:email:%s", s.sec.HMACString(email))
}
func (s *UserStore) usersSetKey() string { return "users:all" }

func (s *UserStore) Create(ctx context.Context, u *model.User) error {
	id := ulid.Make().String()
	now := time.Now().UTC()
	u.ID = id
	u.CreatedAt = now
	u.UpdatedAt = now
	b, _ := json.Marshal(u)
	enc, err := s.sec.Encrypt(b)
	if err != nil {
		return err
	}

	pipe := s.r.TxPipeline()
	pipe.Set(ctx, s.userKey(id), enc, 0)
	pipe.Set(ctx, s.emailKey(u.Email), id, 0)
	pipe.SAdd(ctx, s.usersSetKey(), id)
	_, err = pipe.Exec(ctx)
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
	val, err := s.r.Get(ctx, s.userKey(id)).Result()
	if err != nil {
		return nil, err
	}
	plaintext, err := s.sec.DecryptString(val)
	if err != nil {
		return nil, err
	}
	var u model.User
	if err := json.Unmarshal(plaintext, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) Update(ctx context.Context, u *model.User) error {
	u.UpdatedAt = time.Now().UTC()
	b, _ := json.Marshal(u)
	enc, err := s.sec.Encrypt(b)
	if err != nil {
		return err
	}
	return s.r.Set(ctx, s.userKey(u.ID), enc, 0).Err()
}

func (s *UserStore) ListAllIDs(ctx context.Context) ([]string, error) {
	return s.r.SMembers(ctx, s.usersSetKey()).Result()
}
