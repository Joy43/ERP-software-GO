package supplier

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateSupplierRequest struct {
	OfficeID                      uint    `json:"office_id" binding:"required"`
	Code                          string  `json:"code"`
	Name                          string  `json:"name" binding:"required"`
	VendorName                    string  `json:"vendor_name"`
	MobileNo                      string  `json:"mobile_no" binding:"required"`
	Email                         string  `json:"email"`
	CreditDays                    int     `json:"credit_days"`
	PartnerGroupID               *uint   `json:"partner_group_id"`
	PartnerSubGroupID            *uint   `json:"partner_sub_group_id"`
	Tolerance                     float64 `json:"tolerance"`
	BinNumber                     string  `json:"bin_number"`
	TCSPercentage                 float64 `json:"tcs_percentage"`
	VATRegNoCentral               string  `json:"vat_reg_no_central"`
	TINNumber                     string  `json:"tin_number"`
	TradeLicenseNumber            string  `json:"trade_license_number"`
	IsForeignSupplier             bool    `json:"is_foreign_supplier"`
	VDSApplicable                 bool    `json:"vds_applicable"`
	IsAnonymous                   bool    `json:"is_anonymous"`
	IsVATAccountingNotApplicable  bool    `json:"is_vat_accounting_not_applicable"`
	BillingContactPerson          string  `json:"billing_contact_person"`
	BillingContactNo              string  `json:"billing_contact_no"`
	DistrictID                    *uint   `json:"district_id"`
	ThanaID                       *uint   `json:"thana_id"`
	BillingAddress                string  `json:"billing_address"`
	TaxBracketID                  *uint   `json:"tax_bracket_id"`
	Remarks                       string  `json:"remarks"`
	Attachment                    string  `json:"attachment"`
	Address                       string  `json:"address"`
	GPS                           string  `json:"gps"`
	ContactPerson                 string  `json:"contact_person"`
	VATRegNo                      string  `json:"vat_reg_no"`
}

type UpdateSupplierRequest struct {
	OfficeID                      uint    `json:"office_id"`
	Code                          string  `json:"code"`
	Name                          string  `json:"name"`
	VendorName                    string  `json:"vendor_name"`
	MobileNo                      string  `json:"mobile_no"`
	Email                         string  `json:"email"`
	CreditDays                    int     `json:"credit_days"`
	PartnerGroupID               *uint   `json:"partner_group_id"`
	PartnerSubGroupID            *uint   `json:"partner_sub_group_id"`
	Tolerance                     float64 `json:"tolerance"`
	BinNumber                     string  `json:"bin_number"`
	TCSPercentage                 float64 `json:"tcs_percentage"`
	VATRegNoCentral               string  `json:"vat_reg_no_central"`
	TINNumber                     string  `json:"tin_number"`
	TradeLicenseNumber            string  `json:"trade_license_number"`
	IsForeignSupplier             bool    `json:"is_foreign_supplier"`
	VDSApplicable                 bool    `json:"vds_applicable"`
	IsAnonymous                   bool    `json:"is_anonymous"`
	IsVATAccountingNotApplicable  bool    `json:"is_vat_accounting_not_applicable"`
	BillingContactPerson          string  `json:"billing_contact_person"`
	BillingContactNo              string  `json:"billing_contact_no"`
	DistrictID                    *uint   `json:"district_id"`
	ThanaID                       *uint   `json:"thana_id"`
	BillingAddress                string  `json:"billing_address"`
	TaxBracketID                  *uint   `json:"tax_bracket_id"`
	Remarks                       string  `json:"remarks"`
	Attachment                    string  `json:"attachment"`
	Address                       string  `json:"address"`
	GPS                           string  `json:"gps"`
	ContactPerson                 string  `json:"contact_person"`
	VATRegNo                      string  `json:"vat_reg_no"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateSupplierRequest true "supplier payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/suppliers [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateSupplierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	s := &Supplier{
		OfficeID:                      req.OfficeID,
		Code:                          req.Code,
		Name:                          req.Name,
		VendorName:                    req.VendorName,
		MobileNo:                      req.MobileNo,
		Email:                         req.Email,
		CreditDays:                    req.CreditDays,
		PartnerGroupID:               req.PartnerGroupID,
		PartnerSubGroupID:            req.PartnerSubGroupID,
		Tolerance:                     req.Tolerance,
		BinNumber:                     req.BinNumber,
		TCSPercentage:                 req.TCSPercentage,
		VATRegNoCentral:               req.VATRegNoCentral,
		TINNumber:                     req.TINNumber,
		TradeLicenseNumber:            req.TradeLicenseNumber,
		IsForeignSupplier:             req.IsForeignSupplier,
		VDSApplicable:                 req.VDSApplicable,
		IsAnonymous:                   req.IsAnonymous,
		IsVATAccountingNotApplicable: req.IsVATAccountingNotApplicable,
		BillingContactPerson:          req.BillingContactPerson,
		BillingContactNo:              req.BillingContactNo,
		DistrictID:                    req.DistrictID,
		ThanaID:                       req.ThanaID,
		BillingAddress:                req.BillingAddress,
		TaxBracketID:                  req.TaxBracketID,
		Remarks:                       req.Remarks,
		Attachment:                    req.Attachment,
		Address:                       req.Address,
		GPS:                           req.GPS,
		ContactPerson:                 req.ContactPerson,
		VATRegNo:                      req.VATRegNo,
	}

	if err := h.service.Create(ctx.Request.Context(), s); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create supplier", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "supplier created successfully", s)
}

// List godoc
// @Summary List all suppliers
// @Tags Suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /partners/suppliers [get]
func (h *Handler) List(ctx *gin.Context) {
	suppliers, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch suppliers", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "suppliers fetched successfully", suppliers)
}

// Get godoc
// @Summary Get supplier by ID
// @Tags Suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Supplier ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/suppliers/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	supplier, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "supplier not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "supplier fetched successfully", supplier)
}

// Update godoc
// @Summary Update supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Supplier ID"
// @Param payload body UpdateSupplierRequest true "supplier payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/suppliers/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateSupplierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	s := &Supplier{
		ID:                            uint(id),
		OfficeID:                      req.OfficeID,
		Code:                          req.Code,
		Name:                          req.Name,
		VendorName:                    req.VendorName,
		MobileNo:                      req.MobileNo,
		Email:                         req.Email,
		CreditDays:                    req.CreditDays,
		PartnerGroupID:               req.PartnerGroupID,
		PartnerSubGroupID:            req.PartnerSubGroupID,
		Tolerance:                     req.Tolerance,
		BinNumber:                     req.BinNumber,
		TCSPercentage:                 req.TCSPercentage,
		VATRegNoCentral:               req.VATRegNoCentral,
		TINNumber:                     req.TINNumber,
		TradeLicenseNumber:            req.TradeLicenseNumber,
		IsForeignSupplier:             req.IsForeignSupplier,
		VDSApplicable:                 req.VDSApplicable,
		IsAnonymous:                   req.IsAnonymous,
		IsVATAccountingNotApplicable: req.IsVATAccountingNotApplicable,
		BillingContactPerson:          req.BillingContactPerson,
		BillingContactNo:              req.BillingContactNo,
		DistrictID:                    req.DistrictID,
		ThanaID:                       req.ThanaID,
		BillingAddress:                req.BillingAddress,
		TaxBracketID:                  req.TaxBracketID,
		Remarks:                       req.Remarks,
		Attachment:                    req.Attachment,
		Address:                       req.Address,
		GPS:                           req.GPS,
		ContactPerson:                 req.ContactPerson,
		VATRegNo:                      req.VATRegNo,
	}

	if err := h.service.Update(ctx.Request.Context(), s); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update supplier", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "supplier updated successfully", s)
}

// Delete godoc
// @Summary Delete supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Supplier ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/suppliers/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete supplier", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "supplier deleted successfully", nil)
}
