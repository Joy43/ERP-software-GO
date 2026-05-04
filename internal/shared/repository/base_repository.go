package repository
import (
	"context"
	"gorm.io/gorm"
)

// ------- BaseRepository defines standard CRUD operations for any entity type -------
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindAll(ctx context.Context) ([]T, error)
	FindByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *baseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Model(entity).Updates(entity).Error
}

func (r *baseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}
