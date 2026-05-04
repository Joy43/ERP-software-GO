package wallet

import (
	"context"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[Wallet]
}

type repo struct {
	repository.BaseRepository[Wallet]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[Wallet](db),
		db:             db,
	}
}

func (r *repo) Update(ctx context.Context, wallet *Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}
