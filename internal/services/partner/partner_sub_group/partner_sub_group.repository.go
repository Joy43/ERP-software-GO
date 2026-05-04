package partner_sub_group

import (
	"context"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[PartnerSubGroup]
	FindByGroupID(ctx context.Context, groupID uint) ([]PartnerSubGroup, error)
}

type repositoryImpl struct {
	repository.BaseRepository[PartnerSubGroup]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[PartnerSubGroup](db),
		db:             db,
	}
}

func (r *repositoryImpl) FindAll(ctx context.Context) ([]PartnerSubGroup, error) {
	var subGroups []PartnerSubGroup
	if err := r.db.WithContext(ctx).Preload("PartnerGroup").Find(&subGroups).Error; err != nil {
		return nil, err
	}
	return subGroups, nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id uint) (*PartnerSubGroup, error) {
	var subGroup PartnerSubGroup
	if err := r.db.WithContext(ctx).Preload("PartnerGroup").First(&subGroup, id).Error; err != nil {
		return nil, err
	}
	return &subGroup, nil
}

func (r *repositoryImpl) FindByGroupID(ctx context.Context, groupID uint) ([]PartnerSubGroup, error) {
	var subGroups []PartnerSubGroup
	if err := r.db.WithContext(ctx).Preload("PartnerGroup").Where("partner_group_id = ?", groupID).Find(&subGroups).Error; err != nil {
		return nil, err
	}
	return subGroups, nil
}
