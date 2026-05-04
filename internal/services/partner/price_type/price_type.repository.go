package price_type

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, pt *PriceType) error
	Update(ctx context.Context, pt *PriceType) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*PriceType, error)
	FindAll(ctx context.Context) ([]PriceType, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, pt *PriceType) error {
	return r.db.WithContext(ctx).Create(pt).Error
}

func (r *repository) Update(ctx context.Context, pt *PriceType) error {
	return r.db.WithContext(ctx).Save(pt).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&PriceType{}, id).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*PriceType, error) {
	var pt PriceType
	err := r.db.WithContext(ctx).First(&pt, id).Error
	return &pt, err
}

func (r *repository) FindAll(ctx context.Context) ([]PriceType, error) {
	var pts []PriceType
	err := r.db.WithContext(ctx).Order("id desc").Find(&pts).Error
	return pts, err
}
