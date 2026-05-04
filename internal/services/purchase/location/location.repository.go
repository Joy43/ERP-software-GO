package location

import (
	"errors"
	"gorm.io/gorm"
)

type Repository interface {
	Create(loc *Location) error
	GetByID(id uint) (*Location, error)
	GetByCode(code string) (*Location, error)
	ListAll() ([]Location, error)
	Update(loc *Location) error
	Delete(id uint) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

//
// =============================
// CREATE
// =============================
func (r *repo) Create(loc *Location) error {
	return r.db.Create(loc).Error
}

//
// =============================
// GET BY ID
// =============================
func (r *repo) GetByID(id uint) (*Location, error) {
	var loc Location

	err := r.db.Preload("Office").Preload("Manager").Preload("Parent").Preload("Children").First(&loc, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("location not found")
		}
		return nil, err
	}

	return &loc, nil
}

//
// =============================
// GET BY CODE (for validation)
// =============================
func (r *repo) GetByCode(code string) (*Location, error) {
	var loc Location

	err := r.db.Where("code = ?", code).First(&loc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}

	return &loc, nil
}

//
// =============================
// GET ALL
// =============================
func (r *repo) ListAll() ([]Location, error) {
	var list []Location
	err := r.db.Preload("Office").Preload("Manager").Preload("Parent").Preload("Children").Find(&list).Error
	return list, err
}



//
// =============================
// UPDATE
// =============================
func (r *repo) Update(loc *Location) error {
	return r.db.Save(loc).Error
}

//
// =============================
// DELETE (SOFT DELETE)
// =============================
func (r *repo) Delete(id uint) error {
	return r.db.Delete(&Location{}, id).Error
}