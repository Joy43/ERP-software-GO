package requisitions

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"gorm.io/gorm"
)

// POCreator is a function type to avoid import cycle with purchaseorder package
type POCreator func(requisition *Requisition, createdByID *uint) error

type Service struct {
	repo      *Repository
	poCreator POCreator
}

// NewService creates a new requisition service
func NewService(repo *Repository, poCreator POCreator) *Service {
	return &Service{repo: repo, poCreator: poCreator}
}

// ========================================
//--------- GET ALL------------
// ========================================

func (s *Service) GetAllRequisitions(filter ListRequisitionRequest) (*PaginatedRequisitionResponse, error) {
	total, err := s.repo.Count(filter)
	if err != nil {
		return nil, err
	}

	requisitions, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))

	responses := make([]RequisitionResponse, 0, len(requisitions))
	for _, req := range requisitions {
		responses = append(responses, toRequisitionResponse(&req))
	}

	return &PaginatedRequisitionResponse{
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
// ------GET BY ID------------
// ========================================

func (s *Service) GetRequisitionByID(id uint) (*RequisitionResponse, error) {
	req, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	resp := toRequisitionResponse(req)
	return &resp, nil
}

// ========================================
//--------- CREATE-------------
// ========================================

func (s *Service) CreateRequisition(req *CreateRequisitionRequest, createdByID *uint) (*RequisitionResponse, error) {
	// ----- Type-specific party validation -----
	switch RequisitionType(req.RequisitionType) {
	case EmployeeRequisition:
		// auto-fill employee_id from logged-in user if not provided
		if req.EmployeeID == nil || *req.EmployeeID == 0 {
			if createdByID == nil {
				return nil, errors.New("employee_id is required for EMPLOYEE requisition")
			}
			req.EmployeeID = createdByID
		}
	case DepartmentRequisition:
		if req.DepartmentID == nil || *req.DepartmentID == 0 {
			return nil, errors.New("department_id is required for DEPARTMENT requisition")
		}
	case ProjectRequisition:
		if req.ProjectID == nil || *req.ProjectID == 0 {
			return nil, errors.New("project_id is required for PROJECT requisition")
		}
	}

	// ----- Parse expected date -----
	expectedDate, err := time.Parse("2006-01-02", req.ExpectedDate)
	if err != nil {
		return nil, errors.New("invalid expected_date format, use YYYY-MM-DD")
	}

	// ----- Generate requisition number -----
	reqNumber, err := s.repo.GenerateRequisitionNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate requisition number: %w", err)
	}

	// ----- Build items -----
	items := make([]RequisitionItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, RequisitionItem{
			ItemID:          it.ItemID,
			RequestQuantity: it.RequestQuantity,
			CurrentStock:    it.CurrentStock,
			LastCost:        it.LastCost,
			AverageCost:     it.AverageCost,
			Description:     it.Description,
			ItemTypeID:      it.ItemTypeID,
			CategoryID:      it.CategoryID,
			SubCategoryID:   it.SubCategoryID,
			MinorCategoryID: it.MinorCategoryID,
		})
	}

	// ----- Build model -----
	requisition := &Requisition{
		RequisitionNumber: reqNumber,
		RequisitionType:   RequisitionType(req.RequisitionType),
		Status:            StatusPending,
		ExpectedDate:      expectedDate,
		CreatedDate:       time.Now(),
		Remarks:           req.Remarks,
		Description:       req.Description,
		EmployeeID:        req.EmployeeID,
		DepartmentID:      req.DepartmentID,
		ProjectID:         req.ProjectID,
		BuyerID:           req.BuyerID,
		OfficeID:          req.OfficeID,
		LocationID:        req.LocationID,
		InventoryTypeID:   req.InventoryTypeID,
		SupplierID:        req.SupplierID,
		CreatedByID:       createdByID,
		Items:             items,
	}

	created, err := s.repo.Create(requisition)
	if err != nil {
		return nil, parseDatabaseError(err)
	}

	_ = s.repo.AddStatusHistory(&RequisitionStatusHistory{
		RequisitionID: created.ID,
		UserID:        createdByID,
		ActionType:    ActionCreated,
		FromStatus:    "",
		ToStatus:      string(StatusPending),
	})

	resp := toRequisitionResponse(created)
	return &resp, nil
}

// ========================================
// UPDATE (Pending only)
// ========================================

func (s *Service) UpdateRequisition(id uint, req *UpdateRequisitionRequest, updatedByID *uint) (*RequisitionResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, handleNotFound(err)
	}

	// Update basic fields if provided
	if req.ExpectedDate != nil {
		expectedDate, err := time.Parse("2006-01-02", *req.ExpectedDate)
		if err != nil {
			return nil, errors.New("invalid expected_date format, use YYYY-MM-DD")
		}
		existing.ExpectedDate = expectedDate
	}
	if req.BuyerID != nil {
		existing.BuyerID = req.BuyerID
	}
	if req.OfficeID != nil {
		existing.OfficeID = req.OfficeID
	}
	if req.LocationID != nil {
		existing.LocationID = req.LocationID
	}
	if req.InventoryTypeID != nil {
		existing.InventoryTypeID = req.InventoryTypeID
	}
	if req.SupplierID != nil {
		existing.SupplierID = req.SupplierID
	}
	if req.Remarks != nil {
		existing.Remarks = req.Remarks
	}
	if req.Description != nil {
		existing.Description = req.Description
	}

	// Only update items if they are explicitly provided in the request
	if len(req.Items) > 0 {
		// Check if any items have actual content (not just empty objects)
		hasValidItems := false
		var itemsToValidate []uint
		
		for _, it := range req.Items {
			// Check if this item has any meaningful content
			if it.ItemID != nil && *it.ItemID > 0 {
				hasValidItems = true
				itemsToValidate = append(itemsToValidate, *it.ItemID)
			} else if it.RequestQuantity != nil || it.CurrentStock != nil || it.LastCost != nil || 
					it.AverageCost != nil || it.Description != nil || it.ItemTypeID != nil || 
					it.CategoryID != nil || it.SubCategoryID != nil || it.MinorCategoryID != nil {
				// Item has content but no item_id - this is invalid
				return nil, fmt.Errorf("item_id is required when providing item details")
			}
		}

		// If no valid items provided, skip items update (keep existing items)
		if !hasValidItems {
			// Skip items update - existing items will remain unchanged
		} else {
			// Validate provided item IDs exist in database
			if len(itemsToValidate) > 0 {
				existingItemIDs, err := s.repo.GetExistingItemIDs(itemsToValidate)
				if err != nil {
					return nil, fmt.Errorf("failed to validate item IDs: %w", err)
				}

				if len(existingItemIDs) != len(itemsToValidate) {
					existingMap := make(map[uint]bool)
					for _, id := range existingItemIDs {
						existingMap[id] = true
					}
					var missingIDs []uint
					for _, id := range itemsToValidate {
						if !existingMap[id] {
							missingIDs = append(missingIDs, id)
						}
					}
					return nil, fmt.Errorf("item IDs do not exist: %v", missingIDs)
				}
			}

			// --------- Start Transaction for items update ---------
			tx := s.repo.Begin()

			// --------- Delete old items and create new ones ---------
			if err := tx.DeleteItems(id); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to remove existing items: %w", err)
			}

			// --------- Prepare new items ---------
			newItems := make([]RequisitionItem, 0, len(req.Items))

			for _, it := range req.Items {
				// Only process items that have item_id
				if it.ItemID != nil && *it.ItemID > 0 {
					requestQuantity := float64(1) // default quantity
					if it.RequestQuantity != nil && *it.RequestQuantity > 0 {
						requestQuantity = *it.RequestQuantity
					}

					newItems = append(newItems, RequisitionItem{
						RequisitionID:   id,
						ItemID:          *it.ItemID,
						RequestQuantity: requestQuantity,
						CurrentStock:    it.CurrentStock,
						LastCost:        it.LastCost,
						AverageCost:     it.AverageCost,
						Description:     it.Description,
						ItemTypeID:      it.ItemTypeID,
						CategoryID:      it.CategoryID,
						SubCategoryID:   it.SubCategoryID,
						MinorCategoryID: it.MinorCategoryID,
					})
				}
			}

			// --------- Insert new items ---------
			if len(newItems) > 0 {
				if err := tx.CreateItems(newItems); err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("failed to create items: %w", err)
				}
			}

			// --------- Commit transaction ---------
			if err := tx.Commit(); err != nil {
				return nil, fmt.Errorf("failed to commit transaction: %w", err)
			}

			// Update in memory for response
			existing.Items = newItems
		}
	}
	// If items array is not provided or has no valid items, existing items remain unchanged

	updated, err := s.repo.Update(existing)
	if err != nil {
		return nil, parseDatabaseError(err)
	}

	_ = s.repo.AddStatusHistory(&RequisitionStatusHistory{
		RequisitionID: id,
		UserID:        updatedByID,
		ActionType:    ActionUpdated,
		FromStatus:    string(existing.Status),
		ToStatus:      string(existing.Status),
	})

	resp := toRequisitionResponse(updated)
	return &resp, nil
}

// ========================================
// DELETE (PENDING / CANCELLED only)
// ========================================

func (s *Service) DeleteRequisition(id uint) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return handleNotFound(err)
	}

	if existing.Status != StatusPending && existing.Status != StatusCancelled {
		return fmt.Errorf("requisition cannot be deleted in '%s' status; only PENDING or CANCELLED requisitions can be deleted", existing.Status)
	}

	return s.repo.Delete(id)
}

// ========================================
// UPDATE STATUS (with transition validation)
// ========================================

func (s *Service) UpdateRequisitionStatus(id uint, req *UpdateRequisitionStatusRequest, changedByID *uint) (*RequisitionResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, handleNotFound(err)
	}

	nextStatus := RequisitionStatus(req.Status)

	// ----- Validate transition -----
	if !existing.Status.IsValidTransition(nextStatus) {
		return nil, fmt.Errorf(
			"invalid status transition: '%s' → '%s'. Allowed: %s",
			existing.Status,
			nextStatus,
			allowedTransitionsText(existing.Status),
		)
	}

	fromStatus := string(existing.Status)

	existing.Status = nextStatus
	existing.UpdatedByID = changedByID

	if nextStatus == StatusRejected {
		if req.RejectionReason == nil || strings.TrimSpace(*req.RejectionReason) == "" {
			return nil, errors.New("rejection_reason is required when rejecting a requisition")
		}
		existing.RejectionReason = req.RejectionReason
	}

	// Finance approval: optionally update approved quantities per item

	updated, err := s.repo.Update(existing)
	if err != nil {
		return nil, parseDatabaseError(err)
	}

	// ----- Record status history -----
	_ = s.repo.AddStatusHistory(&RequisitionStatusHistory{
		RequisitionID: id,
		UserID:        changedByID,
		ActionType:    ActionStatusChanged,
		FromStatus:    fromStatus,
		ToStatus:      string(nextStatus),
		Remarks:       req.Remarks,
	})

	// ----- Auto-create PO when status transitions to ORDERED -----
	if nextStatus == StatusOrdered && s.poCreator != nil {
		_ = s.poCreator(updated, changedByID)
	}

	resp := toRequisitionResponse(updated)
	return &resp, nil
}

// ========================================
// GET ALL HISTORY (Value = Quantity × Unit Price)
// ========================================

func (s *Service) GetAllHistory(filter ListHistoryRequest) (*PaginatedHistoryResponse, error) {
	history, total, err := s.repo.GetAllHistory(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))

	result := make([]RequisitionStatusHistoryResponse, 0, len(history))
	for _, h := range history {
		item := RequisitionStatusHistoryResponse{
			ID:            h.ID,
			RequisitionID: h.RequisitionID,
			ActionType:    h.ActionType,
			FromStatus:    h.FromStatus,
			ToStatus:      h.ToStatus,
			Remarks:       h.Remarks,
			CreatedAt:     h.CreatedAt,
		}
		if h.User != nil {
			item.User = &UserDTO{ID: h.User.ID, Name: h.User.Name, Email: h.User.Email}
		}
		result = append(result, item)
	}

	return &PaginatedHistoryResponse{
		Data: result,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			Limit:      filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetRequisitionSummary() (*RequisitionSummaryResponse, error) {
	return s.repo.GetRequisitionSummary()
}

func (s *Service) GetRequisitionHistory(id uint) ([]RequisitionStatusHistoryResponse, error) {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return nil, handleNotFound(err)
	}

	history, err := s.repo.GetStatusHistory(id)
	if err != nil {
		return nil, err
	}

	result := make([]RequisitionStatusHistoryResponse, 0, len(history))
	for _, h := range history {
		item := RequisitionStatusHistoryResponse{
			ID:         h.ID,
			ActionType: h.ActionType,
			FromStatus: h.FromStatus,
			ToStatus:   h.ToStatus,
			Remarks:    h.Remarks,
			CreatedAt:  h.CreatedAt,
		}
		if h.User != nil {
			item.User = &UserDTO{ID: h.User.ID, Name: h.User.Name, Email: h.User.Email}
		}
		result = append(result, item)
	}
	return result, nil
}

// ========================================
// HELPERS
// ========================================

func handleNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	return err
}

func allowedTransitionsText(current RequisitionStatus) string {
	allowed, ok := ValidStatusTransitions[current]
	if !ok || len(allowed) == 0 {
		return "none (terminal state)"
	}
	parts := make([]string, len(allowed))
	for i, s := range allowed {
		parts[i] = string(s)
	}
	return strings.Join(parts, ", ")
}

func parseDatabaseError(err error) error {
	msg := err.Error()

	if errors.Is(err, gorm.ErrInvalidData) {
		return errors.New("invalid data provided")
	}

	if strings.Contains(msg, "foreign key") || strings.Contains(msg, "FOREIGN KEY") {
		fkMap := map[string]string{
			"employee":       "employee_id does not exist",
			"buyer":          "buyer_id does not exist",
			"department":     "department_id does not exist",
			"project":        "project_id does not exist",
			"office":         "office_id does not exist",
			"location":       "location_id does not exist",
			"item":           "item_id does not exist",
			"category":       "category_id does not exist",
			"sub_category":   "sub_category_id does not exist",
			"minor_category": "minor_category_id does not exist",
			"supplier":       "supplier_id does not exist",
			"inventory_type": "inventory_type_id does not exist",
			"item_type":      "item_type_id does not exist",
		}
		for key, message := range fkMap {
			if strings.Contains(strings.ToLower(msg), key) {
				return errors.New(message)
			}
		}
		return fmt.Errorf("referenced entity does not exist: %w", err)
	}

	if strings.Contains(msg, "duplicate") || strings.Contains(msg, "Duplicate") {
		return errors.New("record already exists or violates unique constraint")
	}

	return fmt.Errorf("database error: %w", err)
}

// ========================================
// DTO MAPPERS
// ========================================

func toRequisitionResponse(req *Requisition) RequisitionResponse {
	resp := RequisitionResponse{
		ID:                req.ID,
		RequisitionNumber: req.RequisitionNumber,
		RequisitionType:   req.RequisitionType,
		Status:            req.Status,
		RejectionReason:   req.RejectionReason,
		ExpectedDate:      req.ExpectedDate,
		CreatedDate:       req.CreatedDate,
		Remarks:           req.Remarks,
		Description:       req.Description,
		EmployeeID:        req.EmployeeID,
		DepartmentID:      req.DepartmentID,
		ProjectID:         req.ProjectID,
		BuyerID:           req.BuyerID,
		OfficeID:          req.OfficeID,
		LocationID:        req.LocationID,
		InventoryTypeID:   req.InventoryTypeID,
		SupplierID:        req.SupplierID,
		CreatedByID:       req.CreatedByID,
		UpdatedByID:       req.UpdatedByID,
		CreatedAt:         req.CreatedAt,
		UpdatedAt:         req.UpdatedAt,
	}

	// ----- Relations -----
	if req.Employee != nil {
		resp.Employee = &UserDTO{
			ID:         req.Employee.ID,
			Name:       req.Employee.Name,
			Email:      req.Employee.Email,
			Mobile:     req.Employee.Mobile,
			EmployeeID: req.Employee.EmployeeID,
		}
	}
	if req.Buyer != nil {
		resp.Buyer = &UserDTO{
			ID:         req.Buyer.ID,
			Name:       req.Buyer.Name,
			Email:      req.Buyer.Email,
			Mobile:     req.Buyer.Mobile,
			EmployeeID: req.Buyer.EmployeeID,
		}
	}
	if req.Department != nil {
		resp.Department = &DepartmentDTO{ID: req.Department.ID, Name: req.Department.Name}
	}
	if req.Project != nil {
		resp.Project = &ProjectDTO{ID: req.Project.ID, ProjectName: req.Project.ProjectName}
	}
	if req.Office != nil {
		resp.Office = &OfficeDTO{ID: req.Office.ID, Name: req.Office.Name}
	}
	if req.Location != nil {
		resp.Location = &LocationDTO{ID: req.Location.ID, Name: req.Location.Name, Type: string(req.Location.Type)}
	}
	if req.InventoryType != nil {
		resp.InventoryType = &InventoryTypeDTO{ID: req.InventoryType.ID, Name: req.InventoryType.TypeName}
	}
	if req.Supplier != nil {
		resp.Supplier = &SupplierDTO{ID: req.Supplier.ID, Name: req.Supplier.Name}
	}
	if req.CreatedBy != nil {
		resp.CreatedBy = &UserDTO{
			ID:         req.CreatedBy.ID,
			Name:       req.CreatedBy.Name,
			Email:      req.CreatedBy.Email,
			EmployeeID: req.CreatedBy.EmployeeID,
		}
	}
	if req.UpdatedBy != nil {
		resp.UpdatedBy = &UserDTO{
			ID:         req.UpdatedBy.ID,
			Name:       req.UpdatedBy.Name,
			Email:      req.UpdatedBy.Email,
			EmployeeID: req.UpdatedBy.EmployeeID,
		}
	}

	// ----- Items -----
	resp.Items = make([]RequisitionItemResponse, 0, len(req.Items))
	for _, it := range req.Items {
		itemResp := RequisitionItemResponse{
			ID:               it.ID,
			RequisitionID:    it.RequisitionID,
			ItemID:           it.ItemID,
			RequestQuantity:  it.RequestQuantity,
			ApprovedQuantity: it.ApprovedQuantity,
			CurrentStock:     it.CurrentStock,
			LastCost:         it.LastCost,
			AverageCost:      it.AverageCost,
			Description:      it.Description,
			CreatedAt:        it.CreatedAt,
		}
		if it.Item.ID > 0 {
			itemResp.Item = &ItemDTO{
				ID:      it.Item.ID,
				Name:    it.Item.Name,
				SKU:     it.Item.SKU,
				Barcode: it.Item.Barcode,
			}
		}
		if it.ItemType != nil {
			itemResp.ItemType = &ItemTypeDTO{ID: it.ItemType.ID, Name: it.ItemType.Name}
		}
		if it.Category != nil {
			itemResp.Category = &CategoryDTO{ID: it.Category.ID, Name: it.Category.Name}
		}
		if it.SubCategory != nil {
			itemResp.SubCategory = &SubCategoryDTO{ID: it.SubCategory.ID, Name: it.SubCategory.Name}
		}
		if it.MinorCategory != nil {
			itemResp.MinorCategory = &MinorCategoryDTO{ID: it.MinorCategory.ID, Name: it.MinorCategory.Name}
		}
		resp.Items = append(resp.Items, itemResp)
	}

	// ----- Status History -----
	resp.StatusHistory = make([]RequisitionStatusHistoryResponse, 0, len(req.StatusHistory))
	for _, h := range req.StatusHistory {
		histResp := RequisitionStatusHistoryResponse{
			ID:         h.ID,
			ActionType: h.ActionType,
			FromStatus: h.FromStatus,
			ToStatus:   h.ToStatus,
			Remarks:    h.Remarks,
			CreatedAt:  h.CreatedAt,
		}
		if h.User != nil {
			histResp.User = &UserDTO{
				ID:    h.User.ID,
				Name:  h.User.Name,
				Email: h.User.Email,
			}
		}
		resp.StatusHistory = append(resp.StatusHistory, histResp)
	}

	return resp
}