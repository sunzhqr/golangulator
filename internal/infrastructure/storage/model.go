package storage

import "time"

type HistoryEntryModel struct {
	ID         uint `gorm:"primaryKey"`
	UserID     int64
	Expression string
	Result     float64
	CreatedAt  time.Time
}

func (HistoryEntryModel) TableName() string {
	return "history_entries"
}
