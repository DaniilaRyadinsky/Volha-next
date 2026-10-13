package handlers

import (
	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/gateway/grpc/auth"
	"github.com/autumnterror/volha-backend/internal/gateway/grpc/products"
	"github.com/autumnterror/volha-backend/internal/gateway/redis"
)

type Apis struct {
	apiProduct *products.Client
	apiAuth    *auth.Client
	rds        *redis.Client
	cfg        *config.Config
}

func New(
	apiProduct *products.Client,
	apiAuth *auth.Client,
	rds *redis.Client,
	cfg *config.Config,
) *Apis {
	return &Apis{
		rds:        rds,
		apiProduct: apiProduct,
		apiAuth:    apiAuth,
		cfg:        cfg,
	}
}
