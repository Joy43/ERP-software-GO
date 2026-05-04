package permission

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Permission, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(ctx context.Context) ([]Permission, error) {
	var permissions []Permission
	if err := r.db.WithContext(ctx).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
