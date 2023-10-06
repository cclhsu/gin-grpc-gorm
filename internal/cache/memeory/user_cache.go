package memory

import (
	"context"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type UserMemoryCacheInterface interface {
}

type UserMemoryCache struct {
	ctx			 context.Context
	logger		 *logrus.Logger
	cacheManager *cache.Cache
}

// UserNewMemoryCache creates a new instance of UserMemoryCache.
// It takes a logrus logger and a context as parameters.
func UserNewMemoryCache(ctx context.Context, logger *logrus.Logger, cacheManager *cache.Cache) *UserMemoryCache {
	return &UserMemoryCache{
		ctx:		  ctx,
		logger:		  logger,
		cacheManager: cacheManager,
	}
}
