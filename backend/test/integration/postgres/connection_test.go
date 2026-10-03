//go:build integration

package integration

import (
	"testing"

	"github.com/autumnterror/volha-backend/internal/product-service/config"
	"github.com/autumnterror/volha-backend/internal/product-service/infra/psql"
	"github.com/stretchr/testify/assert"
)

func TestConnectDisconnect(t *testing.T) {
	d, err := psql.NewConnect(config.Test())
	assert.NoError(t, err)
	assert.NoError(t, d.Disconnect())
}
