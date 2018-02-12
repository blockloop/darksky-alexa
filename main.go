package main

import (
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/tea"
	"github.com/caarlos0/env"
	"github.com/garyburd/redigo/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var (
	config = struct {
		DarkskyToken string `env:"DARKSKY_TOKEN,required"`
		RedisURL     string `env:"REDIS_URL"`
		RedisMaxIDLE int    `env:"REDIS_MAX_IDLE" envDefault:"5"`
		Port         int    `env:"PORT", envDefault:"3000"`
	}{}
)

func init() {
	if config.RedisURL == "" {
		// try DOKKU env
		config.RedisURL = os.Getenv("DOKKU_REDIS_DARKSKY_PORT")
	}
}

func main() {
	if err := env.Parse(&config); err != nil {
		log.WithError(err).Fatal("configuration failure")
	}

	tea.Responder = render.JSON

	cache := initCache(config.RedisURL)
	dsapi := initDarksky(config.DarkskyToken)

	mux := chi.NewMux()
	mux.Use(
		middleware.RealIP,
		middleware.RequestID,
		middleware.Timeout(time.Second*10),
		middleware.Logger,
		middleware.Recoverer,
	)

	mux.Post("/darksky", handler(cache, dsapi))
	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"pong": ""}`))
	})

	addr := ":3000"
	log.WithField("addr", addr).Info("server started")
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.WithError(err).Fatalf("shutting down")
	}
	log.Info("shutting down")
}

func initCache(url string) Cache {
	if url == "" {
		return &NopCache{}
	}

	pool := redis.NewPool(func() (redis.Conn, error) {
		return redis.DialURL(config.RedisURL, redis.DialConnectTimeout(5*time.Second))
	}, 5)
	return &RedisCache{pool: pool}
}

func initDarksky(token string) *darksky.API {
	return darksky.New(token)
}
