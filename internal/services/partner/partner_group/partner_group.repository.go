package partner_group

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[PartnerGroup]
}

type repositoryImpl struct {
	repository.BaseRepository[PartnerGroup]
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[PartnerGroup](db),
	}
}
