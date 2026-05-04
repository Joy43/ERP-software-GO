package inventorytypes

import (
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

//---------------  CREATE ----------------
func (s *Service) Create(req CreateInventoryTypeRequest) error {
	
	existing, _ := s.repo.FindByCode(req.TypeCode)
	if existing != nil && existing.ID != 0 {
		return errors.New("type_code already exists")
	}

	inv := InventoryType{
		TypeCode:    req.TypeCode,
		TypeName:    req.TypeName,
		Description: req.Description,
		IsActive:    true,
	}

	if req.IsActive != nil {
		inv.IsActive = *req.IsActive
	}

	return s.repo.Create(&inv)
}

// -------- GET ALL ----------
func (s *Service) GetAll() ([]InventoryType, error) {
	return s.repo.FindAll()
}



// -------------UPDATE----------------
func (s *Service) Update(id int64, req UpdateInventoryTypeRequest) error {
	inv, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("inventory type not found")
	}

	if req.TypeCode != nil {
		existing, err := s.repo.FindByCode(*req.TypeCode)
		if err == nil && existing.ID != inv.ID {
			return errors.New("type_code already exists")
		}
		inv.TypeCode = *req.TypeCode
	}

	if req.TypeName != nil {
		inv.TypeName = *req.TypeName
	}

	if req.Description != nil {
		inv.Description = *req.Description
	}

	if req.IsActive != nil {
		inv.IsActive = *req.IsActive
	}

	return s.repo.Update(inv)
}

//--------------- DELETE (soft delete)-----------------
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}