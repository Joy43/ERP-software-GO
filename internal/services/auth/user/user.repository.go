package user

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, value *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByIDWithRoles(ctx context.Context, id uint) (*User, error)
	FindByEmployeeID(ctx context.Context, employeeID string) (*User, error)
	UpdateUser(ctx context.Context, value *User) error
	DeleteUser(ctx context.Context, id uint) error
	FindAllWithRoles(ctx context.Context) ([]User, error)
	FindAllWithRolesPaginated(ctx context.Context, page, limit int) ([]User, int64, error)
	CreateRefreshToken(ctx context.Context, value *RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uint) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (repository *repository) CreateUser(ctx context.Context, value *User) error {
	if err := repository.db.WithContext(ctx).Create(value).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (repository *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var value User
	if err := repository.db.WithContext(ctx).Where("email = ?", email).First(&value).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

func (repository *repository) FindByID(ctx context.Context, id uint) (*User, error) {
	var value User
	if err := repository.db.WithContext(ctx).First(&value, id).Error; err != nil {
		return nil, err
	}
	return &value, nil
}
func (repository *repository) FindByIDWithRoles(ctx context.Context, id uint) (*User, error) {
	var value User
	if err := repository.db.WithContext(ctx).Preload("Roles.Permissions").First(&value, id).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

func (repository *repository) FindByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	var value User
	if err := repository.db.WithContext(ctx).Where("employee_id = ?", employeeID).First(&value).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

func (repository *repository) UpdateUser(ctx context.Context, value *User) error {
	// Use Save to update all fields including associations if provided, 
	// or Updates for partial updates. Service will handle the logic.
	if err := repository.db.WithContext(ctx).Save(value).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	// Update roles separately since many2many needs explicit handling if using Save on a partial struct
	if len(value.Roles) > 0 {
		if err := repository.db.WithContext(ctx).Model(value).Association("Roles").Replace(value.Roles); err != nil {
			return fmt.Errorf("update user roles: %w", err)
		}
	}
	return nil
}

func (repository *repository) DeleteUser(ctx context.Context, id uint) error {
	return repository.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Revoke all refresh tokens first
		if err := tx.Where("user_id = ?", id).Delete(&RefreshToken{}).Error; err != nil {
			return err
		}
		// Delete user_roles associations
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", id).Error; err != nil {
			return err
		}
		// Delete user
		if err := tx.Delete(&User{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (repository *repository) FindAllWithRoles(ctx context.Context) ([]User, error) {
	var users []User
	if err := repository.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Order("id ASC").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("find all users with roles: %w", err)
	}
	return users, nil
}

// FindAllWithRolesPaginated retrieves users with pagination only
func (repository *repository) FindAllWithRolesPaginated(ctx context.Context, page, limit int) ([]User, int64, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Build query
	query := repository.db.WithContext(ctx)

	// Count total records
	var total int64
	if err := query.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Apply pagination (no sorting)
	var users []User
	if err := query.
		Preload("Roles.Permissions").
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("find all users with roles paginated: %w", err)
	}

	return users, total, nil
}

func (repository *repository) CreateRefreshToken(ctx context.Context, value *RefreshToken) error {
	if err := repository.db.WithContext(ctx).Create(value).Error; err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (repository *repository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	var value RefreshToken
	if err := repository.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now()).
		First(&value).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

func (repository *repository) RevokeRefreshToken(ctx context.Context, id uint) error {
	now := time.Now()
	if err := repository.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", &now).Error; err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (repository *repository) RevokeAllUserRefreshTokens(ctx context.Context, userID uint) error {
	now := time.Now()
	if err := repository.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Update("revoked_at", &now).Error; err != nil {
		return fmt.Errorf("revoke all refresh tokens: %w", err)
	}
	return nil
}
