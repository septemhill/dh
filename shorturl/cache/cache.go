package cache

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type shortURLCache struct {
	client *redis.Client
}

type CacheOption func(*redis.Options)

func Password(passwd string) CacheOption {
	return func(o *redis.Options) {
		o.Password = passwd
	}
}

func DB(db int) CacheOption {
	return func(o *redis.Options) {
		o.DB = db
	}
}

func NewCache(addr string, opts ...CacheOption) *shortURLCache {
	ro := &redis.Options{
		Addr: addr,
	}
	for _, opt := range opts {
		opt(ro)
	}
	client := redis.NewClient(ro)
	return &shortURLCache{
		client: client,
	}
}

// StoreShortURL stores the key, url pair, and expire after an hour.
func (c *shortURLCache) StoreShortURL(key, url string) error {
	_, err := c.client.SetNX(shortURLKeyPrefix+key, url, time.Hour).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *shortURLCache) LookupShortURLByKey(key string) (string, error) {
	url, err := c.client.Get(shortURLKeyPrefix + key).Result()
	if err != nil {
		return "", nil
	}
	return url, nil
}

func (c *shortURLCache) StoreShortURLKeyOffset(offset int64) error {
	_, err := c.client.Set(shortURLKeyOffset, offset, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *shortURLCache) RetrieveShortURLKeyOffset() (int64, error) {
	s, err := c.client.Get(shortURLKeyOffset).Result()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return 0, err
		}

		// Create offset key if not exist, and set to 0.
		c.StoreShortURLKeyOffset(0)
		s = "0"
	}

	off, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return off, nil
}

var _ ShortURLCache = (*shortURLCache)(nil)
