package responsibility_transfer

import (
	"context"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[ResponsibilityTransfer]
	FindAllWithUsers(ctx context.Context, filter TransferFilter) ([]ResponsibilityTransfer, error)
	FindByIDWithUsers(ctx context.Context, id uint) (*ResponsibilityTransfer, error)
}

type repo struct {
	repository.BaseRepository[ResponsibilityTransfer]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[ResponsibilityTransfer](db),
		db:             db,
	}
}

func (r *repo) FindAllWithUsers(ctx context.Context, filter TransferFilter) ([]ResponsibilityTransfer, error) {
	var results []ResponsibilityTransfer
	query := r.db.WithContext(ctx)

	if filter.FromUserID != 0 {
		query = query.Where("from_user_id = ?", filter.FromUserID)
	}
	if filter.ToUserID != 0 {
		query = query.Where("to_user_id = ?", filter.ToUserID)
	}
	if filter.FromDate != nil {
		query = query.Where("from_date >= ?", filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("to_date <= ?", filter.ToDate)
	}

	err := query.
		Preload("FromUser").
		Preload("ToUser").
		Preload("ApprovedBy").
		Order("id DESC").
		Find(&results).Error
	return results, err
}

func (r *repo) FindByIDWithUsers(ctx context.Context, id uint) (*ResponsibilityTransfer, error) {
	var result ResponsibilityTransfer
	err := r.db.WithContext(ctx).
		Preload("FromUser").
		Preload("ToUser").
		Preload("ApprovedBy").
		First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
