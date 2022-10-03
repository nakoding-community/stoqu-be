package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type DurationUnit time.Duration

type Client struct {
	client      *redis.Client
	durationCfg struct {
		duration int
		unit     string
	}
}

func NewClient(c *redis.Client, duration int, du string) *Client {
	client := &Client{
		client: c,
	}
	client.durationCfg.duration = duration
	client.durationCfg.unit = du
	return client
}

func (c *Client) Del(ctx context.Context, key string) error {
	_, err := c.client.Del(ctx, key).Result()
	return err
}

func (c *Client) SetIfNotExist(ctx context.Context, key string, value interface{}) (bool, error) {
	duration := c.durationCfg.duration
	unit := c.durationCfg.unit
	expiration := time.Duration(duration) * time.Duration(StringToLimitUnit(unit))
	return c.client.SetNX(ctx, key, value, expiration).Result()
}
