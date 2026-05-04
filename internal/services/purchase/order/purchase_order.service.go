package purchaseorder

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/requisitions"
	"gorm.io/gorm"
)

type GRNCreator func(po *PurchaseOrder, createdByID *uint) error

type Service struct {
	repo       *Repository
	reqRepo    *requisitions.Repository
	grnCreator GRNCreator
}

func NewService(repo *Repository, reqRepo *requisitions.Repository, grnCreator GRNCreator) *Service {
	return &Service{repo: repo, reqRepo: reqRepo, grnCreator: grnCreator}
}

func (s *Service) GetAll(f ListPORequest) (*PaginatedPOResponse, error) {
	total, err := s.repo.Count(f)
	if err != nil {
		return nil, err
	}
	pos, err := s.repo.FindAll(f)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(f.PageSize)))
	data := make([]POResponse, 0, len(pos))
	for _, po := range pos {
		data = append(data, toPOResponse(&po))
	}
	return &PaginatedPOResponse{
		Data: data,
		Pagination: POPaginationMeta{
			Page: f.Page, Limit: f.PageSize,
			Total: total, TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetByID(id uint) (*POResponse, error) {
	po, err := s.repo.FindByID(id)
	if err != nil {
		return nil, handleNotFound(err)
	}
	resp := toPOResponse(po)
	return &resp, nil
}

func (s *Service) GetOrderedRequisitions() ([]OrderedRequisitionRow, error) {
	return s.repo.GetOrderedRequisitions()
}

func (s *Service) Create(req *CreatePORequest, createdByID *uint) (*POResponse, error) {
	var autoItems []CreatePOItemRequest

	if req.OrderType == POTypeRequisitionBased {
		// ---- Validate requisition ----
		if req.RequisitionID == nil || *req.RequisitionID == 0 {
			return nil, errors.New("requisition_id is required for REQUISITION_BASED order")
		}
		reqRecord, err := s.reqRepo.FindByID(*req.RequisitionID)
		if err != nil {
			return nil, errors.New("requisition not found")
		}
		if reqRecord.Status != requisitions.StatusApproved {
			return nil, fmt.Errorf("requisition must be APPROVED to create a PO, current status: %s", reqRecord.Status)
		}
		if len(reqRecord.Items) == 0 {
			return nil, errors.New("requisition has no items")
		}
		// ---- Auto-populate items from approved requisition items ----
		for _, ri := range reqRecord.Items {
			qty := ri.RequestQuantity
			if ri.ApprovedQuantity != nil && *ri.ApprovedQuantity > 0 {
				qty = *ri.ApprovedQuantity
			}
			unitPrice := float64(0)
			if ri.LastCost != nil {
				unitPrice = *ri.LastCost
			} else if ri.AverageCost != nil {
				unitPrice = *ri.AverageCost
			}
			riID := ri.ID
			autoItems = append(autoItems, CreatePOItemRequest{
				ItemID:            ri.ItemID,
				RequisitionItemID: &riID,
				OrderQuantity:     qty,
				UnitPrice:         unitPrice,
			})
		}
		req.Items = autoItems
	} else {
		// ---- DIRECT / CONTRACT require manual items ----
		if len(req.Items) == 0 {
			return nil, errors.New("items are required for DIRECT and CONTRACT orders")
		}
	}

	poDate, err := time.Parse("2006-01-02", req.PODate)
	if err != nil {
		return nil, errors.New("invalid po_date format, use YYYY-MM-DD")
	}

	var deliveryDate *time.Time
	if req.DeliveryDate != nil && *req.DeliveryDate != "" {
		d, err := time.Parse("2006-01-02", *req.DeliveryDate)
		if err != nil {
			return nil, errors.New("invalid delivery_date format, use YYYY-MM-DD")
		}
		deliveryDate = &d
	}

	poNumber, err := s.repo.GeneratePONumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PO number: %w", err)
	}

	poItems, subtotal, vatTotal, discountTotal := buildPOItems(req.Items)

	po := &PurchaseOrder{
		PONumber:        poNumber,
		PODate:          poDate,
		DeliveryDate:    deliveryDate,
		RequisitionID:   req.RequisitionID,
		OrderType:       req.OrderType,
		OfficeID:        req.OfficeID,
		LocationID:      req.LocationID,
		SupplierID:      req.SupplierID,
		PaymentTerms:    req.PaymentTerms,
		GeneralRemarks:  req.GeneralRemarks,
		ShippingAddress: req.ShippingAddress,
		Status:          POStatusPending,
		CreatedByID:     createdByID,
		Subtotal:        subtotal,
		VatAmount:       vatTotal,
		DiscountAmount:  discountTotal,
		TotalAmount:     subtotal + vatTotal - discountTotal,
		Items:           poItems,
	}

	created, err := s.repo.Create(po)
	if err != nil {
		return nil, parseDatabaseError(err)
	}

	//-------------- Mark requisition as ORDERED --------
	if req.OrderType == POTypeRequisitionBased && req.RequisitionID != nil {
		reqRecord, _ := s.reqRepo.FindByID(*req.RequisitionID)
		if reqRecord != nil {
			reqRecord.Status = requisitions.StatusOrdered
			_, _ = s.reqRepo.Update(reqRecord)
			_ = s.reqRepo.AddStatusHistory(&requisitions.RequisitionStatusHistory{
				RequisitionID: *req.RequisitionID,
				UserID:        createdByID,
				ActionType:    requisitions.ActionStatusChanged,
				FromStatus:    string(requisitions.StatusApproved),
				ToStatus:      string(requisitions.StatusOrdered),
			})
		}
	}

	resp := toPOResponse(created)
	return &resp, nil
}

func (s *Service) Update(id uint, req *UpdatePORequest, updatedByID *uint) (*POResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, handleNotFound(err)
	}
	if existing.Status != POStatusPending {
		return nil, fmt.Errorf("PO can only be edited in PENDING status, current: %s", existing.Status)
	}

	if req.DeliveryDate != nil && *req.DeliveryDate != "" {
		d, err := time.Parse("2006-01-02", *req.DeliveryDate)
		if err != nil {
			return nil, errors.New("invalid delivery_date format")
		}
		existing.DeliveryDate = &d
	}
	if req.PaymentTerms != nil {
		existing.PaymentTerms = req.PaymentTerms
	}
	if req.GeneralRemarks != nil {
		existing.GeneralRemarks = req.GeneralRemarks
	}
	if req.ShippingAddress != nil {
		existing.ShippingAddress = req.ShippingAddress
	}

	if len(req.Items) > 0 {
		if err := s.repo.DeleteItems(id); err != nil {
			return nil, err
		}
		poItems, subtotal, vatTotal, discountTotal := buildPOItems(req.Items)
		for i := range poItems {
			poItems[i].POID = id
		}
		existing.Items = poItems
		existing.Subtotal = subtotal
		existing.VatAmount = vatTotal
		existing.DiscountAmount = discountTotal
		existing.TotalAmount = subtotal + vatTotal - discountTotal
	}

	updated, err := s.repo.Update(existing)
	if err != nil {
		return nil, parseDatabaseError(err)
	}
	resp := toPOResponse(updated)
	return &resp, nil
}

func (s *Service) UpdateStatus(id uint, req *UpdatePOStatusRequest, changedByID *uint) (*POResponse, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, handleNotFound(err)
	}
	if !existing.Status.IsValidTransition(req.Status) {
		return nil, fmt.Errorf("invalid status transition: '%s' → '%s'", existing.Status, req.Status)
	}
	existing.Status = req.Status
	if req.Status == POStatusConfirmed || req.Status == POStatusIssued {
		now := time.Now()
		existing.ApprovedAt = &now
		existing.ApprovedByID = changedByID
	}
	updated, err := s.repo.Update(existing)
	if err != nil {
		return nil, parseDatabaseError(err)
	}

	// Auto-create GRN (AGAINST_PO) when PO is ISSUED
	if req.Status == POStatusIssued && s.grnCreator != nil {
		_ = s.grnCreator(updated, changedByID)
	}

	resp := toPOResponse(updated)
	return &resp, nil
}

func (s *Service) Delete(id uint) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return handleNotFound(err)
	}
	if existing.Status != POStatusPending {
		return fmt.Errorf("only PENDING POs can be deleted, current: %s", existing.Status)
	}
	return s.repo.Delete(id)
}

// =============================================
// HELPERS
// =============================================

func buildPOItems(reqs []CreatePOItemRequest) ([]PurchaseOrderItem, float64, float64, float64) {
	var subtotal, vatTotal, discountTotal float64
	items := make([]PurchaseOrderItem, 0, len(reqs))
	for _, r := range reqs {
		lineTotal := r.OrderQuantity * r.UnitPrice
		vatAmt := lineTotal * r.VatPercentage / 100
		discAmt := lineTotal * r.DiscountPercentage / 100
		total := lineTotal + vatAmt - discAmt

		subtotal += lineTotal
		vatTotal += vatAmt
		discountTotal += discAmt

		items = append(items, PurchaseOrderItem{
			ItemID:             r.ItemID,
			RequisitionItemID:  r.RequisitionItemID,
			OrderQuantity:      r.OrderQuantity,
			UOMID:              r.UOMID,
			UnitPrice:          r.UnitPrice,
			VatPercentage:      r.VatPercentage,
			VatAmount:          vatAmt,
			DiscountPercentage: r.DiscountPercentage,
			DiscountAmount:     discAmt,
			TotalAmount:        total,
			Remarks:            r.Remarks,
		})
	}
	return items, subtotal, vatTotal, discountTotal
}

func toPOResponse(po *PurchaseOrder) POResponse {
	resp := POResponse{
		ID: po.ID, PONumber: po.PONumber, PODate: po.PODate,
		DeliveryDate: po.DeliveryDate, OrderType: po.OrderType, Status: po.Status,
		RequisitionID: po.RequisitionID,
		OfficeID: po.OfficeID, LocationID: po.LocationID, SupplierID: po.SupplierID,
		PaymentTerms: po.PaymentTerms, GeneralRemarks: po.GeneralRemarks,
		ShippingAddress: po.ShippingAddress,
		Subtotal: po.Subtotal, VatAmount: po.VatAmount,
		DiscountAmount: po.DiscountAmount, TotalAmount: po.TotalAmount,
		CreatedByID: po.CreatedByID, ApprovedByID: po.ApprovedByID,
		ApprovedAt: po.ApprovedAt,
		CreatedAt: po.CreatedAt, UpdatedAt: po.UpdatedAt,
	}
	if po.Office != nil {
		resp.Office = &PORefDTO{ID: po.Office.ID, Name: po.Office.Name}
	}
	if po.Location != nil {
		resp.Location = &PORefDTO{ID: po.Location.ID, Name: po.Location.Name}
	}
	if po.Supplier != nil {
		resp.Supplier = &PORefDTO{ID: po.Supplier.ID, Name: po.Supplier.Name}
	}
	resp.Items = make([]POItemResponse, 0, len(po.Items))
	for _, it := range po.Items {
		ir := POItemResponse{
			ID: it.ID, POID: it.POID, ItemID: it.ItemID,
			RequisitionItemID: it.RequisitionItemID,
			OrderQuantity: it.OrderQuantity, ReceivedQuantity: it.ReceivedQuantity,
			UOMID: it.UOMID, UnitPrice: it.UnitPrice,
			VatPercentage: it.VatPercentage, VatAmount: it.VatAmount,
			DiscountPercentage: it.DiscountPercentage, DiscountAmount: it.DiscountAmount,
			TotalAmount: it.TotalAmount, Remarks: it.Remarks, CreatedAt: it.CreatedAt,
		}
		if it.Item.ID > 0 {
			ir.Item = &PORefDTO{ID: it.Item.ID, Name: it.Item.Name}
		}
		resp.Items = append(resp.Items, ir)
	}
	return resp
}

func handleNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	return err
}

func parseDatabaseError(err error) error {
	msg := err.Error()
	if strings.Contains(msg, "foreign key") || strings.Contains(msg, "FOREIGN KEY") {
		return fmt.Errorf("referenced entity does not exist: %w", err)
	}
	if strings.Contains(msg, "duplicate") || strings.Contains(msg, "Duplicate") {
		return errors.New("record already exists or violates unique constraint")
	}
	return fmt.Errorf("database error: %w", err)
}
