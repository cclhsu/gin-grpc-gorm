package memory

import (
	"context"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type TeamMemoryCacheInterface interface {
}

type TeamMemoryCache struct {
	ctx			 context.Context
	logger		 *logrus.Logger
	cacheManager *cache.Cache
}

// TeamNewMemoryCache creates a new instance of TeamMemoryCache.
// It takes a logrus logger and a context as parameters.
func TeamNewMemoryCache(ctx context.Context, logger *logrus.Logger, cacheManager *cache.Cache) *TeamMemoryCache {
	return &TeamMemoryCache{
		ctx:		  ctx,
		logger:		  logger,
		cacheManager: cacheManager,
	}
}
