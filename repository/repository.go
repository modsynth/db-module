package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository provides generic CRUD operations
type Repository[T any] struct {
	db *gorm.DB
}

// New creates a new repository instance
func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// Create creates a new record
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// FindByID finds a record by ID
func (r *Repository[T]) FindByID(ctx context.Context, id interface{}, entity *T) error {
	return r.db.WithContext(ctx).First(entity, id).Error
}

// FindAll finds all records
func (r *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

// Update updates a record
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete deletes a record
func (r *Repository[T]) Delete(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Delete(entity).Error
}

// DeleteByID deletes a record by ID
func (r *Repository[T]) DeleteByID(ctx context.Context, id interface{}) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// Count counts all records
func (r *Repository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).Count(&count).Error
	return count, err
}

// FindWhere finds records matching the condition
func (r *Repository[T]) FindWhere(ctx context.Context, query interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Where(query, args...).Find(&entities).Error
	return entities, err
}

// FirstWhere finds the first record matching the condition
func (r *Repository[T]) FirstWhere(ctx context.Context, entity *T, query interface{}, args ...interface{}) error {
	return r.db.WithContext(ctx).Where(query, args...).First(entity).Error
}

// Paginate returns paginated results
func (r *Repository[T]) Paginate(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64
	var entity T

	// Get total count
	if err := r.db.WithContext(ctx).Model(&entity).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error

	return entities, total, err
}

// Transaction executes operations within a transaction
func (r *Repository[T]) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
