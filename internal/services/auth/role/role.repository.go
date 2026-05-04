package role

import (
	"context"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/permission"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, role *Role) error
	FindAll(ctx context.Context) ([]Role, error)
	FindByID(ctx context.Context, id uint) (*Role, error)
	FindBySlug(ctx context.Context, slug string) (*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uint) error
	ReplacePermissions(ctx context.Context, role *Role, permissions []uint) error
	Transaction(ctx context.Context, fn func(repo Repository) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *repository) FindAll(ctx context.Context) ([]Role, error) {
	var roles []Role

	err := r.db.WithContext(ctx).
		Table("roles").
		Select("roles.*, COUNT(DISTINCT user_roles.user_id) as user_count").
		Joins("LEFT JOIN user_roles ON user_roles.role_id = roles.id").
		Group("roles.id").
		Preload("Permissions").
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
func (r *repository) FindBySlug(ctx context.Context, slug string) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("slug = ?", slug).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) Update(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Role{}, id).Error
}

func (r *repository) ReplacePermissions(ctx context.Context, role *Role, permissionIDs []uint) error {
	var permissions []permission.Permission
	for _, id := range permissionIDs {
		permissions = append(permissions, permission.Permission{ID: id})
	}

	return r.db.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}

func (r *repository) Transaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewRepository(tx))
	})
}
