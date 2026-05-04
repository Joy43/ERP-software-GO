package purchasereturn

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ========================================
// GET ALL
// ========================================

func (s *Service) GetAll(filter ListPurchaseReturnRequest) (*PaginatedPurchaseReturnResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	total, err := s.repo.Count(filter)
	if err != nil {
		return nil, fmt.Errorf("count failed: %w", err)
	}

	records, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, fmt.Errorf("find failed: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))

	responses := make([]PurchaseReturnResponse, 0, len(records))
	for _, r := range records {
		responses = append(responses, toResponse(&r))
	}

	return &PaginatedPurchaseReturnResponse{
		Data: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			Limit:      filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ========================================
//------------- CREATE -------------
// ========================================

func (s *Service) Create(req *CreatePurchaseReturnRequest) (*PurchaseReturnResponse, error) {
	returnNumber, err := s.repo.GenerateReturnNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate return number: %w", err)
	}

	record := &PurchaseReturn{
		ReturnNumber: returnNumber,
		Quantity:     req.Quantity,
		SellingPrice: req.SellingPrice,
		Remarks:      req.Remarks,
		OfficeID:     req.OfficeID,
		LocationID:   req.LocationID,
		SupplierID:   req.SupplierID,
		ItemID:       req.ItemID,
	}

	created, err := s.repo.Create(record)
	if err != nil {
		return nil, parseDatabaseError(err)
	}

	resp := toResponse(created)
	return &resp, nil
}

// ========================================
//----------------  MAPPER ----------------
// ========================================

func toResponse(r *PurchaseReturn) PurchaseReturnResponse {
	resp := PurchaseReturnResponse{
		ID:           r.ID,
		ReturnNumber: r.ReturnNumber,
		Quantity:     r.Quantity,
		SellingPrice: r.SellingPrice,
		Remarks:      r.Remarks,
		OfficeID:     r.OfficeID,
		LocationID:   r.LocationID,
		SupplierID:   r.SupplierID,
		ItemID:       r.ItemID,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}

	if r.Office != nil {
		resp.Office = &OfficeDTO{ID: r.Office.ID, Name: r.Office.Name}
	}
	if r.Location != nil {
		resp.Location = &LocationDTO{ID: r.Location.ID, Name: r.Location.Name, Type: string(r.Location.Type)}
	}
	if r.Supplier != nil {
		resp.Supplier = &SupplierDTO{ID: r.Supplier.ID, Name: r.Supplier.Name}
	}
	if r.Item != nil {
		resp.Item = &ItemDTO{
			ID:      r.Item.ID,
			Name:    r.Item.Name,
			SKU:     r.Item.SKU,
			Barcode: r.Item.Barcode,
		}
	}

	return resp
}

// ========================================
// ERROR HELPER
// ========================================

func parseDatabaseError(err error) error {
	if errors.Is(err, gorm.ErrInvalidData) {
		return errors.New("invalid data provided")
	}

	msg := strings.ToLower(err.Error())

	if strings.Contains(msg, "foreign key") {
		fkMap := map[string]string{
			"office":   "office_id does not exist",
			"location": "location_id does not exist",
			"supplier": "supplier_id does not exist",
			"item":     "item_id does not exist",
		}
		for key, message := range fkMap {
			if strings.Contains(msg, key) {
				return errors.New(message)
			}
		}
		return fmt.Errorf("referenced entity does not exist: %w", err)
	}

	if strings.Contains(msg, "duplicate") {
		return errors.New("record already exists or violates unique constraint")
	}

	return fmt.Errorf("database error: %w", err)
}