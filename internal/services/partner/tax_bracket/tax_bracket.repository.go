package tax_bracket

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[TaxBracket]
}

type repositoryImpl struct {
	repository.BaseRepository[TaxBracket]
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[TaxBracket](db),
	}
}
