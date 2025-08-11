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

type userRecord struct {
	ID                       string    `json:"id"`
	Email                    string    `json:"email"`
	PasswordHash             string    `json:"passwordHash"`
	FullName                 string    `json:"fullName"`
	WeeklyHours              float64   `json:"weeklyHours"`
	DefaultLunchBreakMinutes int       `json:"defaultLunchBreakMinutes"`
	Timezone                 string    `json:"timezone"`
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
}

func toRecord(u *model.User) userRecord {
	return userRecord{
		ID:                       u.ID,
		Email:                    u.Email,
		PasswordHash:             u.PasswordHash,
		FullName:                 u.FullName,
		WeeklyHours:              u.WeeklyHours,
		DefaultLunchBreakMinutes: u.DefaultLunchBreakMinutes,
		Timezone:                 u.Timezone,
		CreatedAt:                u.CreatedAt,
		UpdatedAt:                u.UpdatedAt,
	}
}

func toModel(r userRecord) *model.User {
	return &model.User{
		ID:                       r.ID,
		Email:                    r.Email,
		PasswordHash:             r.PasswordHash,
		FullName:                 r.FullName,
		WeeklyHours:              r.WeeklyHours,
		DefaultLunchBreakMinutes: r.DefaultLunchBreakMinutes,
		Timezone:                 r.Timezone,
		CreatedAt:                r.CreatedAt,
		UpdatedAt:                r.UpdatedAt,
	}
}

func NewUserStore(r redis.UniversalClient, sec *secret.Secret) *UserStore { return &UserStore{r: r, sec: sec} }

func (s *UserStore) userKey(id string) string { return fmt.Sprintf("user:%s", id) }
func (s *UserStore) emailKey(email string) string {
	return fmt.Sprintf("user:email:%s", s.sec.HMACString(email))
}
func (s *UserStore) usersSetKey() string { return "users:all" }

func (s *UserStore) Create(ctx context.Context, u *model.User) error {
	id := ulid.Make().String()
	now := time.Now().UTC()
	u.ID = id
	u.CreatedAt = now
	u.UpdatedAt = now
	rec := toRecord(u)
	b, _ := json.Marshal(rec)
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
	var rec userRecord
	if err := json.Unmarshal(plaintext, &rec); err != nil {
		return nil, err
	}
	return toModel(rec), nil
}

func (s *UserStore) Update(ctx context.Context, u *model.User) error {
	u.UpdatedAt = time.Now().UTC()
	rec := toRecord(u)
	b, _ := json.Marshal(rec)
	enc, err := s.sec.Encrypt(b)
	if err != nil {
		return err
	}
	return s.r.Set(ctx, s.userKey(u.ID), enc, 0).Err()
}

func (s *UserStore) ListAllIDs(ctx context.Context) ([]string, error) {
	return s.r.SMembers(ctx, s.usersSetKey()).Result()
}
