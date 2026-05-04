package designation

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[Designation]
}

type repositoryImpl struct {
	repository.BaseRepository[Designation]
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[Designation](db),
	}
}
