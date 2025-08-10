package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port            int
	RedisURL        string
	JWTSecret       string
	JWTTTLHours     int
	DefaultTimezone string
	AllowOriginsCSV string
}

func Load() Config {
	port := getInt("PORT", 8080)
	redisURL := getString("REDIS_URL", "")
	jwtSecret := getString("JWT_SECRET", "changeme-dev-only")
	jwtTTL := getInt("JWT_TTL_HOURS", 24*14)
	tz := getString("TIMEZONE", "Europe/Paris")
	origins := getString("CORS_ALLOWED_ORIGINS", "*")
	return Config{
		Port:            port,
		RedisURL:        redisURL,
		JWTSecret:       jwtSecret,
		JWTTTLHours:     jwtTTL,
		DefaultTimezone: tz,
		AllowOriginsCSV: origins,
	}
}

func (c Config) HTTPAddr() string { return fmt.Sprintf(":%d", c.Port) }

func getString(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
