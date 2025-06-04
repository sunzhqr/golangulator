package usecase

import (
	"time"

	"github.com/sunzhqr/golangulator/internal/domain"
)

type historyUseCase struct {
	repo domain.HistoryRepository
}

func NewHistoryUseCase(repo domain.HistoryRepository) domain.HistoryUseCase {
	return &historyUseCase{repo: repo}
}

func (h *historyUseCase) SaveEntry(userID int64, expr string, result float64) error {
	entry := &domain.HistoryEntry{
		UserID:     userID,
		Expression: expr,
		Result:     result,
		CreatedAt:  time.Now(),
	}
	return h.repo.Save(entry)
}

func (h *historyUseCase) GetUserHistory(userID int64) ([]domain.HistoryEntry, error) {
	return h.repo.GetByUser(userID)
}

func (h *historyUseCase) ClearUserHistory(userID int64) error {
	return h.repo.ClearByUser(userID)
}
