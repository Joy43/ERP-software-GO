package inventorytypes

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(inv *InventoryType) error {
	return r.db.Create(inv).Error
}

func (r *Repository) FindAll() ([]InventoryType, error) {
	var list []InventoryType
	err := r.db.Find(&list).Error
	return list, err
}

func (r *Repository) FindByID(id int64) (*InventoryType, error) {
	var inv InventoryType
	err := r.db.First(&inv, id).Error
	return &inv, err
}

func (r *Repository) FindByCode(code string) (*InventoryType, error) {
	var inv InventoryType
	err := r.db.Where("type_code = ?", code).First(&inv).Error
	return &inv, err
}

func (r *Repository) Update(inv *InventoryType) error {
	return r.db.Save(inv).Error
}

func (r *Repository) Delete(id int64) error {
	return r.db.Delete(&InventoryType{}, id).Error
}