package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type UserRedisCacheInterface interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd

	// ListUsers() ([]*model.User, error)
	// GetUser(UUID string) *model.User
	// SetUser(user *model.User) error
	// DeleteUser(UUID string) error
	// GetUserByID(ID string) (*model.User, error)
	// GetUserByName(name string) (*model.User, error)
	// GetUserByEmail(email string) (*model.User, error)
}

type UserRedisCache struct {
	ctx	   context.Context
	logger *logrus.Logger
	client *redis.Client
}

func UserNewRedisCache(ctx context.Context, logger *logrus.Logger, client *redis.Client) *UserRedisCache {
	return &UserRedisCache{
		ctx:	ctx,
		logger: logger,
		client: client,
	}
}

func (c *UserRedisCache) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(c.ctx, key, value, expiration)
}

func (c *UserRedisCache) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(c.ctx, key)
}

func (c *UserRedisCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.client.Del(c.ctx, keys...)
}

func (c *UserRedisCache) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.client.Keys(c.ctx, pattern)
}

// func (c *UserRedisCache) ListUsers() ([]*model.User, error) {
//	var users []*model.User
//	keys := c.Keys(c.ctx, "user:*")
//	for _, key := range keys.Val() {
//		user := c.GetUser(key)
//		users = append(users, user)
//	}
//	return users, nil
// }

// func (c *UserRedisCache) GetUserKey(UUID string) string {
//	return "user:" + UUID
// }

// func (c *UserRedisCache) GetUser(UUID string) *model.User {

//	user := &model.User{}
//	val := c.Get(c.ctx, c.GetUserKey(UUID))
//	if val.Err() != nil {
//		return nil
//	}
//	json.Unmarshal([]byte(val.Val()), user)
//	return user
// }

// func (c *UserRedisCache) SetUser(user *model.User) error {
//	userJson, _ := json.Marshal(user)
//	return c.Set(c.GetUserKey(user.UUID), userJson, 0).Err()
// }

// func (c *UserRedisCache) DeleteUser(UUID string) error {
//	return c.Del(c.ctx, c.GetUserKey(UUID)).Err()
// }

// func (c *UserRedisCache) GetUserByID(ID string) (*model.User, error) {
//	return nil, nil
// }

// func (c *UserRedisCache) GetUserByName(name string) (*model.User, error) {
//	return nil, nil
// }

// func (c *UserRedisCache) GetUserByEmail(email string) (*model.User, error) {
//	return nil, nil
// }
