package domain

import "time"

type HistoryEntry struct {
	ID        uint
	UserID    int64
	Expression string
	Result    float64
	CreatedAt time.Time
}

type HistoryRepository interface {
	Save(entry *HistoryEntry) error
	GetByUser(userID int64) ([]HistoryEntry, error)
	ClearByUser(userID int64) error
}

type HistoryUseCase interface {
	SaveEntry(userID int64, expr string, result float64) error
	GetUserHistory(userID int64) ([]HistoryEntry, error)
	ClearUserHistory(userID int64) error
}
