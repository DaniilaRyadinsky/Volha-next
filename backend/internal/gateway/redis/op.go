package redis

import (
	"context"
	"encoding/json"
	"errors"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/model"
	"github.com/autumnterror/volha-backend/internal/platform/format"

	"time"

	"github.com/redis/go-redis/v9"
)

const cacheKey = "dictionaries:all:gateway-v2"
const cacheGenerationKey = "dictionaries:generation"

var ErrCacheMiss = errors.New("cache miss")

func (c *Client) generation(ctx context.Context) (string, error) {
	generation, err := c.Rdb.Get(ctx, cacheGenerationKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "0", nil
		}
		return "", err
	}
	return generation, nil
}

func dictionaryCacheKey(generation string) string {
	return cacheKey + ":v" + generation
}

func dictionaryByCategoryCacheKey(generation, id string) string {
	return cacheKey + ":category:" + id + ":v" + generation
}

func (c *Client) GetDictionaries(ctx context.Context) (*productsRPC.Dictionaries, error) {
	const (
		op = "redis.GetAllDictionaries"
	)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	generation, err := c.generation(ctx)
	if err != nil {
		c.observeCache("dictionaries", "get", "error")
		return nil, format.Error(op, err)
	}
	ds, err := c.Rdb.Get(ctx, dictionaryCacheKey(generation)).Result()
	if err == nil {
		var cached productsRPC.Dictionaries
		if err := json.Unmarshal([]byte(ds), &cached); err != nil {
			c.observeCache("dictionaries", "get", "error")
			return nil, format.Error(op, err)
		}
		c.observeCache("dictionaries", "get", "hit")
		return &cached, nil
	}
	if err == redis.Nil {
		c.observeCache("dictionaries", "get", "miss")
		return nil, format.Error(op, ErrCacheMiss)
	}

	c.observeCache("dictionaries", "get", "error")
	return nil, format.Error(op, err)
}
func (c *Client) SetDictionaries(ctx context.Context, d *productsRPC.Dictionaries) error {
	const (
		op = "redis.SetDictionaries"
	)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	bytes, err := json.Marshal(d)
	if err != nil {
		return format.Error(op, err)
	}

	generation, err := c.generation(ctx)
	if err != nil {
		c.observeCache("dictionaries", "set", "error")
		return format.Error(op, err)
	}
	if err := c.Rdb.Set(ctx, dictionaryCacheKey(generation), bytes, c.ttl).Err(); err != nil {
		c.observeCache("dictionaries", "set", "error")
		return format.Error(op, err)
	}

	c.observeCache("dictionaries", "set", "success")
	return nil
}

func (c *Client) CleanDictionaries(ctx context.Context) error {
	const (
		op = "redis.CleanDictionaries"
	)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := c.Rdb.Incr(ctx, cacheGenerationKey).Err(); err != nil {
		c.observeCache("dictionaries", "invalidate", "error")
		return format.Error(op, err)
	}
	c.observeCache("dictionaries", "invalidate", "success")
	return nil
}

func (c *Client) GetDictionariesByCategory(ctx context.Context, id string) (*model.DictionariesByCategory, error) {
	const (
		op = "redis.GetDictionariesByCategory"
	)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	generation, err := c.generation(ctx)
	if err != nil {
		c.observeCache("dictionaries_by_category", "get", "error")
		return nil, format.Error(op, err)
	}
	ds, err := c.Rdb.Get(ctx, dictionaryByCategoryCacheKey(generation, id)).Result()
	if err == nil {
		var cached model.DictionariesByCategory
		if err := json.Unmarshal([]byte(ds), &cached); err != nil {
			c.observeCache("dictionaries_by_category", "get", "error")
			return nil, format.Error(op, err)
		}
		c.observeCache("dictionaries_by_category", "get", "hit")
		return &cached, nil
	}
	if err == redis.Nil {
		c.observeCache("dictionaries_by_category", "get", "miss")
		return nil, format.Error(op, ErrCacheMiss)
	}

	c.observeCache("dictionaries_by_category", "get", "error")
	return nil, format.Error(op, err)
}
func (c *Client) SetDictionariesByCategory(ctx context.Context, id string, d *model.DictionariesByCategory) error {
	const (
		op = "redis.SetDictionariesByCategory"
	)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	bytes, err := json.Marshal(d)
	if err != nil {
		return format.Error(op, err)
	}

	generation, err := c.generation(ctx)
	if err != nil {
		c.observeCache("dictionaries_by_category", "set", "error")
		return format.Error(op, err)
	}
	if err := c.Rdb.Set(ctx, dictionaryByCategoryCacheKey(generation, id), bytes, c.ttl).Err(); err != nil {
		c.observeCache("dictionaries_by_category", "set", "error")
		return format.Error(op, err)
	}

	c.observeCache("dictionaries_by_category", "set", "success")
	return nil
}
