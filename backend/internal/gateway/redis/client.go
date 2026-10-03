package redis

import (
	"context"
	"time"

	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/redis/go-redis/v9"
)

const serviceName = "volha-gateway"

type Client struct {
	Rdb     *redis.Client
	ttl     time.Duration
	metrics *observability.Metrics
}

func New(cfg *config.Config, metrics *observability.Metrics) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPw,
		DB:       0,
	})
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Warn(ctx, "redis cache unavailable at startup", "operation", "redis.New", "error", err)
	}
	return &Client{Rdb: rdb, ttl: cfg.RedisDictTTL, metrics: metrics}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	const op = "redis.Ping"
	if _, err := c.Rdb.Ping(ctx).Result(); err != nil {
		return format.Error(op, err)
	}
	return nil
}

func (c *Client) observeCache(cache, operation, result string) {
	observability.ObserveRedisCache(c.metrics, serviceName, cache, operation, result)
}

func (c *Client) Close() error {
	const op = "redis.Close"
	if err := c.Rdb.Close(); err != nil {
		return format.Error(op, err)
	}
	return nil
}
