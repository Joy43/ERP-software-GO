package payment_mode

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[PaymentMode]
}

type repositoryImpl struct {
	repository.BaseRepository[PaymentMode]
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[PaymentMode](db),
	}
}
