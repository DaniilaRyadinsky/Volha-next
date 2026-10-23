package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/gateway/model"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/stretchr/testify/require"
)

func TestCacheTTLGenerationAndMalformedValues(t *testing.T) {
	server := miniredis.RunT(t)
	server.RequireAuth("secret")
	client, err := New(&config.Config{RedisAddr: server.Addr(), RedisPw: "secret", RedisDictTTL: time.Second}, observability.NewMetrics("test", "dev", "test"))
	require.NoError(t, err)
	defer client.Close()
	ctx := context.Background()
	require.NoError(t, client.Ping(ctx))
	_, err = client.GetDictionaries(ctx)
	require.ErrorIs(t, err, ErrCacheMiss)
	value := &productsRPC.Dictionaries{MinPrice: 10}
	require.NoError(t, client.SetDictionaries(ctx, value))
	got, err := client.GetDictionaries(ctx)
	require.NoError(t, err)
	require.Equal(t, value.GetMinPrice(), got.GetMinPrice())
	byCategory := &model.DictionariesByCategory{FilteredCategory: "Chairs", MaxPrice: 100}
	require.NoError(t, client.SetDictionariesByCategory(ctx, "cat", byCategory))
	cat, err := client.GetDictionariesByCategory(ctx, "cat")
	require.NoError(t, err)
	require.Equal(t, "Chairs", cat.FilteredCategory)
	require.NoError(t, client.CleanDictionaries(ctx))
	_, err = client.GetDictionaries(ctx)
	require.ErrorIs(t, err, ErrCacheMiss)
	_, err = client.GetDictionariesByCategory(ctx, "cat")
	require.ErrorIs(t, err, ErrCacheMiss)
	require.NoError(t, client.SetDictionaries(ctx, value))
	server.FastForward(2 * time.Second)
	_, err = client.GetDictionaries(ctx)
	require.ErrorIs(t, err, ErrCacheMiss)
	require.NoError(t, server.Set(dictionaryCacheKey("1"), "{broken"))
	_, err = client.GetDictionaries(ctx)
	require.Error(t, err)
	require.NoError(t, server.Set(dictionaryByCategoryCacheKey("1", "cat"), "{broken"))
	_, err = client.GetDictionariesByCategory(ctx, "cat")
	require.Error(t, err)
	require.NoError(t, client.Close())
	require.Error(t, client.Ping(ctx))
	_, err = client.GetDictionaries(ctx)
	require.Error(t, err)
	require.Error(t, client.CleanDictionaries(ctx))
}
