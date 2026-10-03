//go:build integration

package integration

import (
	"context"
	"errors"
	"testing"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/gateway/model"
	redis "github.com/autumnterror/volha-backend/internal/gateway/redis"
	"github.com/stretchr/testify/assert"
)

func TestDictionariesGood(t *testing.T) {
	c, err := redis.New(config.Test(), nil)
	assert.NoError(t, err)

	td := &productsRPC.Dictionaries{
		Brands:     []*productsRPC.Brand{},
		Categories: []*productsRPC.Category{},
		Countries:  []*productsRPC.Country{},
		Materials:  []*productsRPC.Material{},
		Colors:     []*productsRPC.Color{},
		MinPrice:   1,
		MaxPrice:   2,
		MinWidth:   3,
		MaxWidth:   4,
		MinHeight:  5,
		MaxHeight:  6,
		MinDepth:   7,
		MaxDepth:   8,
	}

	assert.NoError(t, c.SetDictionaries(context.Background(), td))

	d, err := c.GetDictionaries(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, d)
	t.Logf("dictionaries cache hit: brands=%d categories=%d", len(d.Brands), len(d.Categories))

	assert.NoError(t, c.CleanDictionaries(context.Background()))

	d, err = c.GetDictionaries(context.Background())
	assert.ErrorIs(t, err, redis.ErrCacheMiss)
	assert.Nil(t, d)
}

func TestDictionariesNoCache(t *testing.T) {
	c, err := redis.New(config.Test(), nil)
	assert.NoError(t, err)

	assert.NoError(t, c.CleanDictionaries(context.Background()))

	d, err := c.GetDictionaries(context.Background())
	assert.ErrorIs(t, err, redis.ErrCacheMiss)
	assert.Nil(t, d)
}

func TestDictionariesByCategoryGood(t *testing.T) {
	c, err := redis.New(config.Test(), nil)
	assert.NoError(t, err)

	id := "category-1"

	td := &model.DictionariesByCategory{
		Brands:     []*productsRPC.Brand{},
		Categories: []*productsRPC.Category{},
		Countries:  []*productsRPC.Country{},
		Materials:  []*productsRPC.Material{},
		Colors:     []*productsRPC.Color{},
		MinPrice:   1,
		MaxPrice:   2,
		MinWidth:   3,
		MaxWidth:   4,
		MinHeight:  5,
		MaxHeight:  6,
		MinDepth:   7,
		MaxDepth:   8,
	}

	assert.NoError(t, c.SetDictionariesByCategory(context.Background(), id, td))

	d, err := c.GetDictionariesByCategory(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, d)
	t.Logf("dictionaries by category cache hit: brands=%d categories=%d", len(d.Brands), len(d.Categories))

	assert.Equal(t, td.MinPrice, d.MinPrice)
	assert.Equal(t, td.MaxPrice, d.MaxPrice)
	assert.Equal(t, td.MinWidth, d.MinWidth)
	assert.Equal(t, td.MaxWidth, d.MaxWidth)
	assert.Equal(t, td.MinHeight, d.MinHeight)
	assert.Equal(t, td.MaxHeight, d.MaxHeight)
	assert.Equal(t, td.MinDepth, d.MinDepth)
	assert.Equal(t, td.MaxDepth, d.MaxDepth)

	assert.NoError(t, c.CleanDictionaries(context.Background()))

	d, err = c.GetDictionariesByCategory(context.Background(), id)
	assert.True(t, errors.Is(err, redis.ErrCacheMiss))
	assert.Nil(t, d)
}

func TestDictionariesByCategoryNoCache(t *testing.T) {
	c, err := redis.New(config.Test(), nil)
	assert.NoError(t, err)

	id := "category-2"

	assert.NoError(t, c.CleanDictionaries(context.Background()))

	d, err := c.GetDictionariesByCategory(context.Background(), id)
	assert.ErrorIs(t, err, redis.ErrCacheMiss)
	assert.Nil(t, d)
}
