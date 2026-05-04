package transaction_history

import (
	"context"
)

type Service interface {
	CreateHistory(ctx context.Context, deletedByID uint, req CreateDeleteHistoryRequest) (*TransactionDeleteHistory, error)
	ListHistory(ctx context.Context, filter HistoryFilter) ([]TransactionDeleteHistory, error)
	DeleteHistory(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateHistory(ctx context.Context, deletedByID uint, req CreateDeleteHistoryRequest) (*TransactionDeleteHistory, error) {
	history := &TransactionDeleteHistory{
		TransactionType: req.TransactionType,
		InvoiceRefNo:    req.InvoiceRefNo,
		Amount:          req.Amount,
		DeletedByID:     &deletedByID,
		ReasonRemarks:   req.ReasonRemarks,
	}

	if err := s.repo.Create(ctx, history); err != nil {
		return nil, err
	}
	// Redetch with preloaded DeletedBy
	var result *TransactionDeleteHistory
	var err error
	// Here we just use FindAllWithFilters with filter on this ID for preloading or a specialized find method
	results, err := s.repo.FindAllWithFilters(ctx, HistoryFilter{TransactionType: history.TransactionType}) // Simple reload
	if err == nil && len(results) > 0 {
		for _, r := range results {
			if r.ID == history.ID {
				result = &r
				break
			}
		}
	}
	if result == nil {
		return history, nil
	}
	return result, nil
}

func (s *service) ListHistory(ctx context.Context, filter HistoryFilter) ([]TransactionDeleteHistory, error) {
	return s.repo.FindAllWithFilters(ctx, filter)
}

func (s *service) DeleteHistory(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
