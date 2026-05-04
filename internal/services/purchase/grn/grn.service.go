package grn

import (
	"errors"
	"fmt"
	"math"
	"time"

	purchaseorder "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/order"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/stock"
	"gorm.io/gorm"
)

type Service struct {
	repo    *Repository
	poRepo  *purchaseorder.Repository
	stockRepo *stock.Repository
}

func NewService(repo *Repository, poRepo *purchaseorder.Repository, stockRepo *stock.Repository) *Service {
	return &Service{repo: repo, poRepo: poRepo, stockRepo: stockRepo}
}

func (s *Service) GetAll(f ListGRNRequest) (*PaginatedGRNResponse, error) {
	total, err := s.repo.Count(f)
	if err != nil {
		return nil, err
	}
	grns, err := s.repo.FindAll(f)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(f.PageSize)))
	data := make([]GRNResponse, 0, len(grns))
	for _, g := range grns {
		data = append(data, toGRNResponse(&g))
	}
	return &PaginatedGRNResponse{
		Data: data,
		Pagination: GRNPaginationMeta{
			Page: f.Page, PageSize: f.PageSize,
			Total: total, TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetByID(id uint) (*GRNResponse, error) {
	g, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	resp := toGRNResponse(g)
	return &resp, nil
}

// Create saves a GRN as PENDING. Call Approve to confirm and update stock.
func (s *Service) Create(req *CreateGRNRequest, createdByID *uint) (*GRNResponse, error) {
	// ---- Validate AGAINST_PO ----
	if req.ReceiveType == GRNAgainstPO {
		if req.POID == nil || *req.POID == 0 {
			return nil, errors.New("po_id is required for AGAINST_PO receive type")
		}
		po, err := s.poRepo.FindByID(*req.POID)
		if err != nil {
			return nil, errors.New("purchase order not found")
		}
		if po.Status == purchaseorder.POStatusCancelled || po.Status == purchaseorder.POStatusFullyReceived {
			return nil, fmt.Errorf("cannot receive against PO in '%s' status", po.Status)
		}
		if err := s.validatePOItems(po, req.Items); err != nil {
			return nil, err
		}
	}

	grnDate, err := time.Parse("2006-01-02", req.GRNDate)
	if err != nil {
		return nil, errors.New("invalid grn_date format, use YYYY-MM-DD")
	}

	var challanDate *time.Time
	if req.ChallanDate != nil && *req.ChallanDate != "" {
		d, err := time.Parse("2006-01-02", *req.ChallanDate)
		if err != nil {
			return nil, errors.New("invalid challan_date format")
		}
		challanDate = &d
	}

	grnNumber, err := s.repo.GenerateGRNNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate GRN number: %w", err)
	}

	grnItems := buildGRNItems(req.Items)

	g := &GoodsReceiptNote{
		GRNNumber:              grnNumber,
		GRNDate:                grnDate,
		ReceiveType:            req.ReceiveType,
		Status:                 GRNStatusPending,
		POID:                   req.POID,
		RequisitionID:          req.RequisitionID,
		OfficeID:               req.OfficeID,
		LocationID:             req.LocationID,
		SupplierID:             req.SupplierID,
		ChallanNo:              req.ChallanNo,
		ChallanDate:            challanDate,
		SalesInvoiceNumber:     req.SalesInvoiceNumber,
		VATChallanNumber:       req.VATChallanNumber,
		DeliveryNumber:         req.DeliveryNumber,
		ShippingAddress:        req.ShippingAddress,
		ShipmentDocumentNumber: req.ShipmentDocumentNumber,
		PaymentMethodID:        req.PaymentMethodID,
		Remarks:                req.Remarks,
		FileID:                 req.FileID,
		CreatedByID:            createdByID,
		ReceivedByID:           req.ReceivedByID,
		Items:                  grnItems,
	}

	created, err := s.repo.Create(g)
	if err != nil {
		return nil, fmt.Errorf("failed to create GRN: %w", err)
	}

	resp := toGRNResponse(created)
	return &resp, nil
}

// Approve confirms a DRAFT GRN:
// 1. Updates location_stocks (weighted average cost)
// 2. Records stock_transactions per item
// 3. Updates purchase_order_items.received_quantity (if AGAINST_PO)
// 4. Recalculates PO status (PARTIALLY_RECEIVED / FULLY_RECEIVED)
// 5. Sets GRN status to CONFIRMED
// All steps run inside a single DB transaction for atomicity.
func (s *Service) Approve(id uint, approvedByID *uint) (*GRNResponse, error) {
	g, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if g.Status != GRNStatusPending {
		return nil, fmt.Errorf("only PENDING GRNs can be approved, current status: %s", g.Status)
	}

	// Re-validate PO items before confirming
	if g.ReceiveType == GRNAgainstPO && g.POID != nil {
		po, err := s.poRepo.FindByID(*g.POID)
		if err != nil {
			return nil, errors.New("purchase order not found")
		}
		if po.Status == purchaseorder.POStatusCancelled || po.Status == purchaseorder.POStatusFullyReceived {
			return nil, fmt.Errorf("cannot approve GRN: PO is in '%s' status", po.Status)
		}
		grnItemReqs := make([]CreateGRNItemRequest, 0, len(g.Items))
		for _, it := range g.Items {
			grnItemReqs = append(grnItemReqs, CreateGRNItemRequest{
				ItemID:           it.ItemID,
				POItemID:         it.POItemID,
				ReceivedQuantity: it.ReceivedQuantity,
			})
		}
		if err := s.validatePOItems(po, grnItemReqs); err != nil {
			return nil, err
		}
	}

	if err := s.confirm(g, approvedByID); err != nil {
		return nil, err
	}

	final, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := toGRNResponse(final)
	return &resp, nil
}

func (s *Service) confirm(g *GoodsReceiptNote, userID *uint) error {
	db := s.repo.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		for _, item := range g.Items {
			txNum, err := s.stockRepo.GenerateTxNumber()
			if err != nil {
				return err
			}

			// Get current stock before update
			currentStock, _ := s.stockRepo.GetStock(item.ItemID, g.LocationID)
			var beforeQty float64
			if currentStock != nil {
				beforeQty = currentStock.Quantity
			}

			// Upsert location_stocks (weighted average cost)
			ls, err := s.stockRepo.UpsertStock(tx, item.ItemID, g.LocationID, item.ReceivedQuantity, item.PurchasePrice)
			if err != nil {
				return fmt.Errorf("failed to update stock for item %d: %w", item.ItemID, err)
			}

			// Record stock transaction
			grnItemID := item.ID
			if err := s.stockRepo.CreateTransaction(tx, &stock.StockTransaction{
				TransactionNumber: txNum,
				TransactionType:   stock.TxGRN,
				ItemID:            item.ItemID,
				LocationID:        g.LocationID,
				QuantityChange:    item.ReceivedQuantity,
				BeforeQuantity:    beforeQty,
				AfterQuantity:     ls.Quantity,
				UnitCost:          item.PurchasePrice,
				ReferenceType:     "GRN",
				ReferenceID:       g.ID,
				GRNItemID:         &grnItemID,
				CreatedByID:       userID,
			}); err != nil {
				return fmt.Errorf("failed to record stock transaction: %w", err)
			}

			// Update PO item received quantity
			if g.ReceiveType == GRNAgainstPO && item.POItemID != nil {
				if err := s.poRepo.UpdateItemReceivedQty(tx, *item.POItemID, item.ReceivedQuantity); err != nil {
					return fmt.Errorf("failed to update PO item received qty: %w", err)
				}
			}
		}

		// Recalculate PO status
		if g.ReceiveType == GRNAgainstPO && g.POID != nil {
			if err := s.poRepo.UpdatePOReceiveStatus(tx, *g.POID); err != nil {
				return fmt.Errorf("failed to update PO status: %w", err)
			}
		}

		// Mark GRN as CONFIRMED
		return s.repo.UpdateStatus(tx, g.ID, GRNStatusConfirmed)
	})
}

// validatePOItems checks that each GRN item references a valid PO item
// and that received qty doesn't exceed remaining (order - already received).
func (s *Service) validatePOItems(po *purchaseorder.PurchaseOrder, items []CreateGRNItemRequest) error {
	poItemMap := make(map[uint]*purchaseorder.PurchaseOrderItem, len(po.Items))
	for i := range po.Items {
		poItemMap[po.Items[i].ID] = &po.Items[i]
	}
	for _, item := range items {
		if item.POItemID == nil {
			return fmt.Errorf("po_item_id is required for item_id %d when receive_type is AGAINST_PO", item.ItemID)
		}
		poItem, ok := poItemMap[*item.POItemID]
		if !ok {
			return fmt.Errorf("po_item_id %d does not belong to PO %d", *item.POItemID, po.ID)
		}
		remaining := poItem.OrderQuantity - poItem.ReceivedQuantity
		if item.ReceivedQuantity > remaining {
			return fmt.Errorf(
				"received_quantity %.2f exceeds remaining %.2f for po_item_id %d",
				item.ReceivedQuantity, remaining, *item.POItemID,
			)
		}
	}
	return nil
}

// =============================================
// HELPERS
// =============================================

func buildGRNItems(reqs []CreateGRNItemRequest) []GRNItem {
	items := make([]GRNItem, 0, len(reqs))
	for _, r := range reqs {
		lineTotal := r.ReceivedQuantity * r.PurchasePrice
		vatAmt := lineTotal * r.VatPercentage / 100
		discAmt := lineTotal * r.DiscountPercentage / 100
		items = append(items, GRNItem{
			ItemID:             r.ItemID,
			POItemID:           r.POItemID,
			ReceivedQuantity:   r.ReceivedQuantity,
			UOMID:              r.UOMID,
			PurchasePrice:      r.PurchasePrice,
			VatPercentage:      r.VatPercentage,
			VatAmount:          vatAmt,
			DiscountPercentage: r.DiscountPercentage,
			DiscountAmount:     discAmt,
			TotalAmount:        lineTotal + vatAmt - discAmt,
			CategoryID:         r.CategoryID,
			SubCategoryID:      r.SubCategoryID,
			MinorCategoryID:    r.MinorCategoryID,
			Remarks:            r.Remarks,
		})
	}
	return items
}

func toGRNResponse(g *GoodsReceiptNote) GRNResponse {
	resp := GRNResponse{
		ID: g.ID, GRNNumber: g.GRNNumber, GRNDate: g.GRNDate,
		ReceiveType: g.ReceiveType, Status: g.Status,
		POID: g.POID, RequisitionID: g.RequisitionID,
		OfficeID: g.OfficeID, LocationID: g.LocationID, SupplierID: g.SupplierID,
		ChallanNo: g.ChallanNo, ChallanDate: g.ChallanDate,
		SalesInvoiceNumber: g.SalesInvoiceNumber, VATChallanNumber: g.VATChallanNumber,
		DeliveryNumber: g.DeliveryNumber, ShippingAddress: g.ShippingAddress,
		ShipmentDocumentNumber: g.ShipmentDocumentNumber,
		PaymentMethodID: g.PaymentMethodID, Remarks: g.Remarks,
		FileID: g.FileID,
		CreatedByID: g.CreatedByID, ReceivedByID: g.ReceivedByID,
		CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
	if g.Office != nil {
		resp.Office = &GRNRefDTO{ID: g.Office.ID, Name: g.Office.Name}
	}
	if g.Location != nil {
		resp.Location = &GRNRefDTO{ID: g.Location.ID, Name: g.Location.Name}
	}
	if g.Supplier != nil {
		resp.Supplier = &GRNRefDTO{ID: g.Supplier.ID, Name: g.Supplier.Name}
	}
	if g.File != nil {
		resp.File = &FileDTO{
			ID:           g.File.ID,
			Name:         g.File.Name,
			OriginalName: g.File.OriginalName,
			MimeType:     g.File.MimeType,
			Size:         g.File.Size,
			StoragePath:  g.File.StoragePath,
			Extension:    g.File.Extension,
			IsPublic:     g.File.IsPublic,
		}
	}
	resp.Items = make([]GRNItemResponse, 0, len(g.Items))
	for _, it := range g.Items {
		ir := GRNItemResponse{
			ID: it.ID, GRNID: it.GRNID, ItemID: it.ItemID, POItemID: it.POItemID,
			ReceivedQuantity: it.ReceivedQuantity, UOMID: it.UOMID,
			PurchasePrice: it.PurchasePrice,
			VatPercentage: it.VatPercentage, VatAmount: it.VatAmount,
			DiscountPercentage: it.DiscountPercentage, DiscountAmount: it.DiscountAmount,
			TotalAmount: it.TotalAmount,
			CategoryID: it.CategoryID, SubCategoryID: it.SubCategoryID,
			MinorCategoryID: it.MinorCategoryID, Remarks: it.Remarks,
			CreatedAt: it.CreatedAt,
		}
		if it.Item.ID > 0 {
			ir.Item = &GRNRefDTO{ID: it.Item.ID, Name: it.Item.Name}
		}
		resp.Items = append(resp.Items, ir)
	}
	return resp
}
