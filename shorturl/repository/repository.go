package repository

import (
	"log"
	"strconv"
	"time"

	//_ "github.com/lib/pq"

	"github.com/septemhill/dh/shorturl/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type shortURLRepository struct {
	db *gorm.DB
}

// Assume we already have a table "short_url_table",
// and created by following SQL command:
//
// CREATE TABLE short_url_table (id SERIAL PRIMARY KEY, key VARCHAR(10), url VARCHAR(512), expired timestamp);

// NewRepository creates a postgres repository for save short url information
func NewRepository(connStr string) *shortURLRepository {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect to postgres database")
	}

	return &shortURLRepository{
		db: db,
	}
}

// StoreURL saves url into short_url_table.
func (repo *shortURLRepository) StoreURL(key, url string, expired time.Time) error {
	return repo.db.Create(models.ShortURL{
		Key:     key,
		URL:     url,
		Expired: expired.UTC(),
	}).Error
}

func (repo *shortURLRepository) StoreLastKey(count int) error {
	return nil
}

// FindURLByID searches specified record by id from short_url_table.
func (repo *shortURLRepository) FindURLByID(id int) (*models.ShortURL, error) {
	f := &models.ShortURL{
		ID: strconv.FormatInt(int64(id), 10),
	}
	if err := repo.db.First(f, f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

// FindURLByID searches specified record by url from short_url_table.
func (repo *shortURLRepository) FindURLByURL(url string, expired time.Time) (*models.ShortURL, error) {
	f := &models.ShortURL{
		URL: url,
	}
	if err := repo.db.Where("url = ? AND expired > ?", url, time.Now().Format(time.RFC3339)).First(f).Error; err != nil {
		return nil, err
	}

	if expired.After(f.Expired.Local()) {
		if err := repo.db.Model(&models.ShortURL{ID: f.ID}).Update("expired", expired.UTC()).Error; err != nil {
			return nil, err
		}
	}

	f.Expired = f.Expired.Local()
	return f, nil
}

// FindURLByKey searches specified record by key from short_url_table.
func (repo *shortURLRepository) FindURLByKey(key string) (*models.ShortURL, error) {
	f := &models.ShortURL{
		Key: key,
	}
	if err := repo.db.First(f, f).Error; err != nil {
		return nil, err
	}
	f.Expired = f.Expired.Local()
	return f, nil
}

var _ ShortURLRepository = (*shortURLRepository)(nil)
