package store

import (
	"crypto/tls"
	"errors"
	"strings"

	"work-tracker/internal/secret"

	"github.com/redis/go-redis/v9"
)

type Stores struct {
	Users       *UserStore
	Sessions    *SessionStore
	TimeLogs    *TimeLogStore
	TimeSheets  *TimeSheetStore
	MonthRecaps *MonthRecapStore
	DaysOff     *DaysOffStore
}

type cryptoDeps struct{ sec *secret.Secret }

func NewStores(c redis.UniversalClient, sec *secret.Secret) Stores {
	return Stores{
		Users:       NewUserStore(c, sec),
		Sessions:    NewSessionStore(c),
		TimeLogs:    NewTimeLogStore(c, sec),
		TimeSheets:  NewTimeSheetStore(c, sec),
		MonthRecaps: NewMonthRecapStore(c, sec),
		DaysOff:     NewDaysOffStore(c, sec),
	}
}

func NewRedisClient(redisURL string) (redis.UniversalClient, error) {
	if redisURL == "" {
		return nil, errors.New("REDIS_URL is required")
	}
	// upstash rediss URL supported by ParseURL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(redisURL, "rediss://") {
		if opt.TLSConfig == nil {
			opt.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
	}
	return redis.NewClient(opt), nil
}
