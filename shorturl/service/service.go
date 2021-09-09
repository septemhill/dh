package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/septemhill/dh/shorturl/cache"
	"github.com/septemhill/dh/shorturl/repository"
	"github.com/septemhill/dh/shorturl/utils"
	"gorm.io/gorm"
)

const baseURL = "http://localhost/"

type service struct {
	repo   repository.ShortURLRepository
	cache  cache.ShortURLCache
	keys   chan string
	offset int64
}

func NewShortURLService(ctx context.Context, repo repository.ShortURLRepository, cache cache.ShortURLCache) *service {
	srv := &service{
		repo:  repo,
		cache: cache,
		keys:  make(chan string, 1000),
	}

	offset, err := srv.cache.RetrieveShortURLKeyOffset()
	if err != nil {
		panic("failed to get latest offset from cache: " + err.Error())
	}

	srv.offset = offset
	return srv
}

func (s *service) generateKey() string {
	key := utils.TransformNumberTo62Digit(s.offset)
	s.offset++
	s.cache.StoreShortURLKeyOffset(s.offset)
	return fmt.Sprintf("%06s", key)
}

// UploadURL creates a short url.
func (s *service) UploadURL(url string, expired time.Time) (string, string, error) {
	sh, err := s.repo.FindURLByURL(url, expired)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", errors.Wrap(err, "failed to query from db")
	}

	// If the URL still not expired,
	// DO NOT generate another one, use previous created.
	if sh != nil && time.Now().Before(sh.Expired) {
		return sh.Key, baseURL + sh.Key, nil
	}

	key := s.generateKey()
	s.repo.StoreURL(key, url, expired)
	s.cache.StoreShortURL(key, url)

	return key, baseURL + key, nil
}

// AccessURL check the short url exist or not.
func (s *service) AccessURL(key string) (string, error) {
	url, err := s.cache.LookupShortURLByKey(key)
	if err != nil {
		return "", err
	}

	// Cache hit, return the URL.
	if url != "" {
		return url, nil
	}

	// Cache miss, search from database.
	sh, err := s.repo.FindURLByKey(key)
	if err != nil {
		return "", err
	}

	// Short URL found, but expired, so return empty url without error.
	if time.Now().After(sh.Expired) {
		return "", nil
	}

	return sh.URL, nil
}

var _ ShortURLService = (*service)(nil)
