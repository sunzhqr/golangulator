package storage

import (
	"github.com/sunzhqr/golangulator/internal/domain"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(entry *domain.HistoryEntry) error {
	model := HistoryEntryModel{
		UserID:     entry.UserID,
		Expression: entry.Expression,
		Result:     entry.Result,
		CreatedAt:  entry.CreatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *PostgresRepo) GetByUser(userID int64) ([]domain.HistoryEntry, error) {
	var models []HistoryEntryModel
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(10).Find(&models).Error
	if err != nil {
		return nil, err
	}
	entries := make([]domain.HistoryEntry, len(models))
	for i, m := range models {
		entries[i] = domain.HistoryEntry{
			ID:         m.ID,
			UserID:     m.UserID,
			Expression: m.Expression,
			Result:     m.Result,
			CreatedAt:  m.CreatedAt,
		}
	}
	return entries, nil
}

func (r *PostgresRepo) ClearByUser(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&HistoryEntryModel{}).Error
}
