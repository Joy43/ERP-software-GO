package transaction_history

import (
	"context"
	"fmt"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[TransactionDeleteHistory]
	FindAllWithFilters(ctx context.Context, filter HistoryFilter) ([]TransactionDeleteHistory, error)
}

type repo struct {
	repository.BaseRepository[TransactionDeleteHistory]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[TransactionDeleteHistory](db),
		db:             db,
	}
}

func (r *repo) FindAllWithFilters(ctx context.Context, filter HistoryFilter) ([]TransactionDeleteHistory, error) {
	var results []TransactionDeleteHistory
	query := r.db.WithContext(ctx).Preload("DeletedBy")

	if filter.TransactionType != "" && filter.TransactionType != "All Types" {
		query = query.Where("transaction_type = ?", filter.TransactionType)
	}

	if filter.DeletedBy != 0 {
		query = query.Where("deleted_by_id = ?", filter.DeletedBy)
	}

	if filter.FromDate != "" {
		query = query.Where("DATE(deleted_at) >= ?", filter.FromDate)
	}

	if filter.ToDate != "" {
		query = query.Where("DATE(deleted_at) <= ?", filter.ToDate)
	}

	err := query.Order("id DESC").Find(&results).Error
	if err != nil {
		fmt.Printf("Error execution query in TransactionDeleteHistory: %v\n", err)
		return nil, err
	}
	return results, nil
}
