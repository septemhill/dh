package models

import "time"

type ShortURL struct {
	ID      string    `gorm:"->" json:"id"`
	Key     string    `gorm:"column:key" json:"key"`
	URL     string    `gorm:"column:url" json:"url"`
	Expired time.Time `gorm:"column:expired" json:"expiredAt"`
}

func (s ShortURL) TableName() string {
	return "short_url_table"
}
