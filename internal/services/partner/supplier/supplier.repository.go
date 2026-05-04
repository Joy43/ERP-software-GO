package supplier

import (
	"context"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[Supplier]
	FindAllWithPreload(ctx context.Context) ([]Supplier, error)
	FindByIDWithPreload(ctx context.Context, id uint) (*Supplier, error)
}

type repositoryImpl struct {
	repository.BaseRepository[Supplier]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[Supplier](db),
		db:             db,
	}
}

func (r *repositoryImpl) FindAllWithPreload(ctx context.Context) ([]Supplier, error) {
	var suppliers []Supplier
	err := r.db.WithContext(ctx).
		Preload("Office").
		Preload("PartnerGroup").
		Preload("PartnerSubGroup").
		Preload("District").
		Preload("Thana").
		Preload("TaxBracket").
		Find(&suppliers).Error
	return suppliers, err
}

func (r *repositoryImpl) FindByIDWithPreload(ctx context.Context, id uint) (*Supplier, error) {
	var s Supplier
	err := r.db.WithContext(ctx).
		Preload("Office").
		Preload("PartnerGroup").
		Preload("PartnerSubGroup").
		Preload("District").
		Preload("Thana").
		Preload("TaxBracket").
		First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
