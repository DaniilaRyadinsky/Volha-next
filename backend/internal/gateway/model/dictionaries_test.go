package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDictionaryCacheJSONUsesCamelCase(t *testing.T) {
	raw, err := json.Marshal(DictionariesByCategory{MinPrice: 1, MaxWidth: 2})
	require.NoError(t, err)
	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got))
	require.Contains(t, got, "minPrice")
	require.Contains(t, got, "maxWidth")
	require.NotContains(t, got, "min_price")
}
