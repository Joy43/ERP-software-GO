package sales_setup

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[SalesSetup]
}

type repositoryImpl struct {
	repository.BaseRepository[SalesSetup]
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[SalesSetup](db),
	}
}
