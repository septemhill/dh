package service

import "time"

type ShortURLService interface {
	UploadURL(url string, expired time.Time) (string, string, error)
	AccessURL(key string) (string, error)
}
