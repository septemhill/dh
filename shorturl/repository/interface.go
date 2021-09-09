package repository

import (
	"time"

	"github.com/septemhill/dh/shorturl/models"
)

type ShortURLRepository interface {
	StoreURL(key string, url string, expired time.Time) error
	StoreLastKey(count int) error
	FindURLByID(id int) (*models.ShortURL, error)
	FindURLByURL(url string, expired time.Time) (*models.ShortURL, error)
	FindURLByKey(key string) (*models.ShortURL, error)
}
