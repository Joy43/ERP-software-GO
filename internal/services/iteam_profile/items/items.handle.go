package items

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

//----------- Handler manages HTTP requests for create item operations -------------
type Handler struct {
	service *Service
	rdb     *redis.Client
}

// NewHandler creates and returns a new handler instance
func NewHandler(service *Service, rdb *redis.Client) *Handler {
	return &Handler{
		service: service,
		rdb:     rdb,
	}
}



// Create godoc
// @Summary Create a new item
// @Description Creates a new item with all required fields
// @Tags Create-Items--------------items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateItemRequest true "item creation payload"
// @Success 201 {object} response.APIResponse{data=ItemResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile-items/items [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	// -----------Create model from request-------------
	item := &Items{
		Name:                req.Name,
		Barcode:             req.Barcode,
		SKU:                 req.SKU,
		CategoryID:          req.CategoryID,
		SubCategoryID:       req.SubCategoryID,
		MinorCategoryID:     req.MinorCategoryID,
		ItemTypeID:          req.ItemTypeID,
		DepartmentID:        req.DepartmentID,
		UomId:               req.UomID,
		SupplierID:          req.SupplierID,
		TagID:               req.TagID,
		CostPrice:           req.CostPrice,
		StandardSalesPrice:  req.StandardSalesPrice,
		CostOnGP:            req.CostOnGP,
		LastCost:            req.LastCost,
		AvgCost:             req.AvgCost,
		PriceIncreasing:     req.PriceIncreasing,
		AltUMO:              req.AltUMO,
		AltUnit:             req.AltUnit,
		MaxDiscount:         req.MaxDiscount,
		ReorderMinimumQty:   req.ReorderMinimumQty,
		CalculateBasePrice:  req.CalculateBasePrice,
		SalesSupplyTypeID:   req.SalesSupplyTypeID,
		SalesTaxSetupID:     req.SalesTaxSetupID,
		SalesSetupID:        req.SalesSetupID,
		FileID:              req.FileID,
		IsChildBarcode:      req.IsChildBarcode,
		AutoBarcode:         req.AutoBarcode,
		CanBeSold:           req.CanBeSold,
		CanBeProduced:       req.CanBeProduced,
		CanBeRented:         req.CanBeRented,
		CanBePurchased:      req.CanBePurchased,
		IsVATRebatable:      req.IsVATRebatable,
		IsNotAllowDecimal:   req.IsNotAllowDecimal,
		IsActive:            req.IsActive,
		IsStyle:             req.IsStyle,
		IsPercentage:        req.IsPercentage,
		IsPharmacy:          req.IsPharmacy,
	}

	// Convert pharmacy request to model if provided
	var pharmacy *ItemPharmacy
	if req.IsPharmacy && req.Pharmacy != nil {
		pharmacy = convertCreatePharmacyRequestToModel(req.Pharmacy)
	}

	if err := h.service.Create(ctx.Request.Context(), item, pharmacy); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to create item", "CREATION_ERROR", err.Error())
		return
	}

	// Invalidate list cache when new item is created


	itemResponse := convertToItemResponse(item)
	response.Success(ctx, "item created successfully", itemResponse)
}

// List godoc
// @Summary List all items with pagination
// @Description Retrieves all items with optional pagination and sorting
// @Tags Create-Items--------------items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order (asc/desc)"
// @Success 200 {object} response.APIResponse{data=ListResponse}
// @Failure 500 {object} response.APIResponse
// @Router /profile-items/items [get]
func (h *Handler) List(ctx *gin.Context) {
	var pagination PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid pagination parameters", "VALIDATION_ERROR", err.Error())
		return
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 20
	}




	items, total, err := h.service.FindAllPaginated(ctx.Request.Context(), pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch items", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))

	itemResponses := make([]ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = convertToItemResponse(&item)
	}

	listResp := ListResponse{
		Items:      itemResponses,
		Total:      total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: totalPages,
	}

	


	response.Success(ctx, "items fetched successfully", listResp)
}

// Get godoc
// @Summary Get a specific item by ID
// @Description Retrieves a single item by its ID
// @Tags Create-Items--------------items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Item ID"
// @Success 200 {object} response.APIResponse{data=ItemResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile-items/items/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid item id format", "INVALID_ID", err.Error())
		return
	}

	itemID := uint(id)



	// --------- Cache miss - query database ----------
	item, err := h.service.FindByID(ctx.Request.Context(), itemID)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "item not found", "NOT_FOUND", err.Error())
		return
	}

	itemResponse := convertToItemResponse(item)

	

	response.Success(ctx, "item fetched successfully", itemResponse)
}

// Update godoc
// @Summary Update an item
// @Description Updates an existing item with provided fields
// @Tags Create-Items--------------items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Item ID"
// @Param payload body UpdateItemRequest true "item update payload"
// @Success 200 {object} response.APIResponse{data=ItemResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile-items/items/{id} [patch]
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid item id format", "INVALID_ID", err.Error())
		return
	}

	var req UpdateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	// --------Fetch existing item--------
	item, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "item not found", "NOT_FOUND", err.Error())
		return
	}


	applyUpdateFields(item, &req)

	if err := h.service.Update(ctx.Request.Context(), item); err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to update item", "UPDATE_ERROR", err.Error())
		return
	}


	itemResponse := convertToItemResponse(item)
	response.Success(ctx, "item updated successfully", itemResponse)
}

// Delete godoc
// @Summary Delete an item
// @Description Soft deletes an item by its ID
// @Tags Create-Items--------------items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Item ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile-items/items/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid item id format", "INVALID_ID", err.Error())
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusNotFound, "failed to delete item", "DELETION_ERROR", err.Error())
		return
	}



	response.Success(ctx, "item deleted successfully", nil)
}


func applyUpdateFields(item *Items, req *UpdateItemRequest) {
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Barcode != nil {
		item.Barcode = *req.Barcode
	}
	if req.SKU != nil {
		item.SKU = *req.SKU
	}
	if req.CategoryID != nil {
		item.CategoryID = *req.CategoryID
	}
	if req.SubCategoryID != nil {
		item.SubCategoryID = *req.SubCategoryID
	}
	if req.MinorCategoryID != nil {
		item.MinorCategoryID = *req.MinorCategoryID
	}
	if req.ItemTypeID != nil {
		item.ItemTypeID = *req.ItemTypeID
	}
	if req.DepartmentID != nil {
		item.DepartmentID = *req.DepartmentID
	}
	if req.UomID != nil {
		item.UomId = *req.UomID
	}
	if req.SupplierID != nil {
		item.SupplierID = *req.SupplierID
	}
	if req.TagID != nil {
		item.TagID = *req.TagID
	}
	if req.CostPrice != nil {
		item.CostPrice = *req.CostPrice
	}
	if req.StandardSalesPrice != nil {
		item.StandardSalesPrice = *req.StandardSalesPrice
	}
	if req.CostOnGP != nil {
		item.CostOnGP = *req.CostOnGP
	}
	if req.LastCost != nil {
		item.LastCost = *req.LastCost
	}
	if req.AvgCost != nil {
		item.AvgCost = *req.AvgCost
	}
	if req.PriceIncreasing != nil {
		item.PriceIncreasing = *req.PriceIncreasing
	}
	if req.AltUMO != nil {
		item.AltUMO = *req.AltUMO
	}
	if req.AltUnit != nil {
		item.AltUnit = *req.AltUnit
	}
	if req.MaxDiscount != nil {
		item.MaxDiscount = *req.MaxDiscount
	}
	if req.ReorderMinimumQty != nil {
		item.ReorderMinimumQty = *req.ReorderMinimumQty
	}
	if req.CalculateBasePrice != nil {
		item.CalculateBasePrice = *req.CalculateBasePrice
	}
	if req.SalesSupplyTypeID != nil {
		item.SalesSupplyTypeID = *req.SalesSupplyTypeID
	}
	if req.SalesTaxSetupID != nil {
		item.SalesTaxSetupID = *req.SalesTaxSetupID
	}
	if req.SalesSetupID != nil {
		item.SalesSetupID = *req.SalesSetupID
	}
	if req.FileID != nil {
		item.FileID = req.FileID
	}

	if req.IsChildBarcode != nil {
		item.IsChildBarcode = *req.IsChildBarcode
	}
	if req.AutoBarcode != nil {
		item.AutoBarcode = *req.AutoBarcode
	}
	if req.CanBeSold != nil {
		item.CanBeSold = *req.CanBeSold
	}
	if req.CanBeProduced != nil {
		item.CanBeProduced = *req.CanBeProduced
	}
	if req.CanBeRented != nil {
		item.CanBeRented = *req.CanBeRented
	}
	if req.CanBePurchased != nil {
		item.CanBePurchased = *req.CanBePurchased
	}
	if req.IsVATRebatable != nil {
		item.IsVATRebatable = *req.IsVATRebatable
	}
	if req.IsNotAllowDecimal != nil {
		item.IsNotAllowDecimal = *req.IsNotAllowDecimal
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if req.IsStyle != nil {
		item.IsStyle = *req.IsStyle
	}
	if req.IsPercentage != nil {
		item.IsPercentage = *req.IsPercentage
	}
}

//------ convertToItemResponse converts a Items model to ItemResponse------
func convertToItemResponse(item *Items) ItemResponse {
	return ItemResponse{
		ID:                 item.ID,
		Name:               item.Name,
		Barcode:            item.Barcode,
		SKU:                item.SKU,
		CategoryID:         item.CategoryID,
		Category:           item.Category,
		SubCategoryID:      item.SubCategoryID,
		SubCategory:        item.SubCategory,
		MinorCategoryID:    item.MinorCategoryID,
		MinorCategory:      item.MinorCategory,
		ItemTypeID:         item.ItemTypeID,
		ItemType:           item.ItemType,
		DepartmentID:       item.DepartmentID,
		Department:         item.Department,
		UomID:              item.UomId,
		Uom:                item.Uom,
		SupplierID:         item.SupplierID,
		Supplier:           item.Supplier,
		TagID:              item.TagID,
		Tag:                item.Tag,
		SalesSupplyTypeID:  item.SalesSupplyTypeID,
		SalesSupplyType:    item.SalesSupplyType,
		SalesTaxSetupID:    item.SalesTaxSetupID,
		SalesTaxSetup:      item.SalesTaxSetup,
		SalesSetupID:       item.SalesSetupID,
		SalesSetup:         item.SalesSetup,
		FileID:             item.FileID,
		File:               item.File,
		CostPrice:          item.CostPrice,
		StandardSalesPrice: item.StandardSalesPrice,
		CostOnGP:           item.CostOnGP,
		LastCost:           item.LastCost,
		AvgCost:            item.AvgCost,
		PriceIncreasing:    item.PriceIncreasing,
		AltUMO:             item.AltUMO,
		AltUnit:            item.AltUnit,
		MaxDiscount:        item.MaxDiscount,
		ReorderMinimumQty:  item.ReorderMinimumQty,
		CalculateBasePrice: item.CalculateBasePrice,
		IsChildBarcode:     item.IsChildBarcode,
		AutoBarcode:        item.AutoBarcode,
		CanBeSold:          item.CanBeSold,
		CanBeProduced:      item.CanBeProduced,
		CanBeRented:        item.CanBeRented,
		CanBePurchased:     item.CanBePurchased,
		IsVATRebatable:     item.IsVATRebatable,
		IsNotAllowDecimal:  item.IsNotAllowDecimal,
		IsActive:           item.IsActive,
		IsStyle:            item.IsStyle,
		IsPercentage:       item.IsPercentage,
		IsPharmacy:         item.IsPharmacy,
		Pharmacy:           convertPharmacyToResponse(item.Pharmacy),
		CreatedAt:          item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// convertPharmacyToResponse converts ItemPharmacy model to ItemPharmacyResponse
func convertPharmacyToResponse(pharmacy *ItemPharmacy) *ItemPharmacyResponse {
	if pharmacy == nil {
		return nil
	}
	return &ItemPharmacyResponse{
		ItemID:                 pharmacy.ItemID,
		GenericName:            pharmacy.GenericName,
		BrandName:              pharmacy.BrandName,
		Strength:               pharmacy.Strength,
		DosageForm:             pharmacy.DosageForm,
		ScheduleType:           pharmacy.ScheduleType,
		IsPrescriptionRequired: pharmacy.IsPrescriptionRequired,
		IsControlledDrug:       pharmacy.IsControlledDrug,
		StorageCondition:       pharmacy.StorageCondition,
		MaxDailyDose:           pharmacy.MaxDailyDose,
		ShelfLifeDays:          pharmacy.ShelfLifeDays,
		ReorderAlertDays:       pharmacy.ReorderAlertDays,
		ManufacturerName:       pharmacy.ManufacturerName,
		DrugRegistrationNo:     pharmacy.DrugRegistrationNo,
		RouteOfAdministration:  pharmacy.RouteOfAdministration,
		TherapeuticClass:       pharmacy.TherapeuticClass,
		CreatedAt:              pharmacy.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:              pharmacy.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// convertCreatePharmacyRequestToModel converts CreateItemPharmacyRequest to ItemPharmacy model
func convertCreatePharmacyRequestToModel(req *CreateItemPharmacyRequest) *ItemPharmacy {
	if req == nil {
		return nil
	}
	return &ItemPharmacy{ 
		GenericName:            req.GenericName,
		BrandName:              req.BrandName,
		Strength:               req.Strength,
		DosageForm:             req.DosageForm,
		ScheduleType:           req.ScheduleType,
		IsPrescriptionRequired: req.IsPrescriptionRequired,
		IsControlledDrug:       req.IsControlledDrug,
		StorageCondition:       req.StorageCondition,
		MaxDailyDose:           req.MaxDailyDose,
		ShelfLifeDays:          req.ShelfLifeDays,
		ReorderAlertDays:       req.ReorderAlertDays,
		ManufacturerName:       req.ManufacturerName,
		DrugRegistrationNo:     req.DrugRegistrationNo,
		RouteOfAdministration:  req.RouteOfAdministration,
		TherapeuticClass:       req.TherapeuticClass,
	}
}