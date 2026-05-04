package sales_representative

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, sr *SalesRepresentative) error
	Update(ctx context.Context, sr *SalesRepresentative) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*SalesRepresentative, error)
	FindAll(ctx context.Context) ([]SalesRepresentative, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, sr *SalesRepresentative) error {
	return r.db.WithContext(ctx).Create(sr).Error
}

func (r *repository) Update(ctx context.Context, sr *SalesRepresentative) error {
	return r.db.WithContext(ctx).Save(sr).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&SalesRepresentative{}, id).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*SalesRepresentative, error) {
	var sr SalesRepresentative
	err := r.db.WithContext(ctx).First(&sr, id).Error
	return &sr, err
}

func (r *repository) FindAll(ctx context.Context) ([]SalesRepresentative, error) {
	var srs []SalesRepresentative
	err := r.db.WithContext(ctx).Order("id desc").Find(&srs).Error
	return srs, err
}
