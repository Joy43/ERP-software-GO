package thana

import (
	"context"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[Thana]
	FindByDistrictID(ctx context.Context, districtID uint) ([]Thana, error)
}

type repositoryImpl struct {
	repository.BaseRepository[Thana]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[Thana](db),
		db:             db,
	}
}

func (r *repositoryImpl) FindAll(ctx context.Context) ([]Thana, error) {
	var thanas []Thana
	if err := r.db.WithContext(ctx).Preload("District").Find(&thanas).Error; err != nil {
		return nil, err
	}
	return thanas, nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id uint) (*Thana, error) {
	var thana Thana
	if err := r.db.WithContext(ctx).Preload("District").First(&thana, id).Error; err != nil {
		return nil, err
	}
	return &thana, nil
}

func (r *repositoryImpl) FindByDistrictID(ctx context.Context, districtID uint) ([]Thana, error) {
	var thanas []Thana
	if err := r.db.WithContext(ctx).Preload("District").Where("district_id = ?", districtID).Find(&thanas).Error; err != nil {
		return nil, err
	}
	return thanas, nil
}
