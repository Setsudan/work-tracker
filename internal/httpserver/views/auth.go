package views

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"work-tracker/internal/auth"
	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/middleware"
	"work-tracker/internal/model"
	"work-tracker/internal/store"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

type registerReq struct {
	Email                    string  `json:"email"`
	Password                 string  `json:"password"`
	FullName                 string  `json:"fullName"`
	WeeklyHours              float64 `json:"weeklyHours"`
	DefaultLunchBreakMinutes int     `json:"defaultLunchBreakMinutes"`
	Timezone                 string  `json:"timezone"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResp struct {
	Token string `json:"token"`
}

func Auth(cfg config.Config, stores store.Stores) chi.Router {
	r := chi.NewRouter()

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		var req registerReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" || req.Password == "" {
			http.Error(w, "email and password required", http.StatusBadRequest)
			return
		}
		if _, err := stores.Users.GetByEmail(r.Context(), req.Email); err == nil {
			http.Error(w, "email in use", http.StatusConflict)
			return
		}
		hash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		u := &model.User{
			Email:                    req.Email,
			PasswordHash:             hash,
			FullName:                 req.FullName,
			WeeklyHours:              req.WeeklyHours,
			DefaultLunchBreakMinutes: req.DefaultLunchBreakMinutes,
			Timezone:                 ifEmpty(req.Timezone, cfg.DefaultTimezone),
		}
		if err := stores.Users.Create(r.Context(), u); err != nil {
			http.Error(w, "could not create user", http.StatusInternalServerError)
			return
		}
		jti := ulid.Make().String()
		token, err := auth.IssueToken(u.ID, cfg.JWTSecret, time.Duration(cfg.JWTTTLHours)*time.Hour, jti)
		if err != nil {
			http.Error(w, "could not issue token", http.StatusInternalServerError)
			return
		}
		_ = stores.Sessions.Create(r.Context(), jti, u.ID, time.Duration(cfg.JWTTTLHours)*time.Hour)
		writeJSON(w, http.StatusOK, tokenResp{Token: token})
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		u, err := stores.Users.GetByEmail(r.Context(), strings.ToLower(strings.TrimSpace(req.Email)))
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		ok, _ := argon2id.ComparePasswordAndHash(req.Password, u.PasswordHash)
		if !ok {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		jti := ulid.Make().String()
		token, err := auth.IssueToken(u.ID, cfg.JWTSecret, time.Duration(cfg.JWTTTLHours)*time.Hour, jti)
		if err != nil {
			http.Error(w, "could not issue token", http.StatusInternalServerError)
			return
		}
		_ = stores.Sessions.Create(r.Context(), jti, u.ID, time.Duration(cfg.JWTTTLHours)*time.Hour)
		writeJSON(w, http.StatusOK, tokenResp{Token: token})
	})

	r.With(middleware.WithAuth(cfg, stores)).Get("/me", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		u, err := stores.Users.GetByID(r.Context(), uid)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		u.PasswordHash = ""
		writeJSON(w, http.StatusOK, u)
	})

	r.With(middleware.WithAuth(cfg, stores)).Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 {
			http.Error(w, "bad token", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ParseToken(parts[1], cfg.JWTSecret)
		if err == nil {
			_ = stores.Sessions.Delete(r.Context(), claims.ID)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return r
}

func ifEmpty(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
