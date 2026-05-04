package sales_supply_type

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[SalesSupplyType]
}

type repositoryImpl struct {
	repository.BaseRepository[SalesSupplyType]
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[SalesSupplyType](db),
	}
}