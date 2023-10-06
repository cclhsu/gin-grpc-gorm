package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type TeamRedisCacheInterface interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd

	// ListTeams() ([]*model.Team, error)
	// GetTeam(UUID string) *model.Team
	// SetTeam(team *model.Team) error
	// DeleteTeam(UUID string) error
	// GetTeamByID(ID string) (*model.Team, error)
	// GetTeamByName(name string) (*model.Team, error)
	// GetTeamByEmail(email string) (*model.Team, error)
}

type TeamRedisCache struct {
	ctx	   context.Context
	logger *logrus.Logger
	client *redis.Client
}

func TeamNewRedisCache(ctx context.Context, logger *logrus.Logger, client *redis.Client) *TeamRedisCache {
	return &TeamRedisCache{
		ctx:	ctx,
		logger: logger,
		client: client,
	}
}

func (c *TeamRedisCache) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(c.ctx, key, value, expiration)
}

func (c *TeamRedisCache) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(c.ctx, key)
}

func (c *TeamRedisCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.client.Del(c.ctx, keys...)
}

func (c *TeamRedisCache) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.client.Keys(c.ctx, pattern)
}

// func (c *TeamRedisCache) ListTeams() ([]*model.Team, error) {
//	var teams []*model.Team
//	keys := c.Keys(c.ctx, "team:*")
//	for _, key := range keys.Val() {
//		team := c.GetTeam(key)
//		teams = append(teams, team)
//	}
//	return teams, nil
// }

// func (c *TeamRedisCache) GetTeamKey(UUID string) string {
//	return "team:" + UUID
// }

// func (c *TeamRedisCache) GetTeam(UUID string) *model.Team {

//	team := &model.Team{}
//	val := c.Get(c.ctx, c.GetTeamKey(UUID))
//	if val.Err() != nil {
//		return nil
//	}
//	json.Unmarshal([]byte(val.Val()), team)
//	return team
// }

// func (c *TeamRedisCache) SetTeam(team *model.Team) error {
//	teamJson, _ := json.Marshal(team)
//	return c.Set(c.GetTeamKey(team.UUID), teamJson, 0).Err()
// }

// func (c *TeamRedisCache) DeleteTeam(UUID string) error {
//	return c.Del(c.ctx, c.GetTeamKey(UUID)).Err()
// }

// func (c *TeamRedisCache) GetTeamByID(ID string) (*model.Team, error) {
//	return nil, nil
// }

// func (c *TeamRedisCache) GetTeamByName(name string) (*model.Team, error) {
//	return nil, nil
// }

// func (c *TeamRedisCache) GetTeamByEmail(email string) (*model.Team, error) {
//	return nil, nil
// }
