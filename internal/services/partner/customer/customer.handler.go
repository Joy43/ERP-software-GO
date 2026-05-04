package customer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/utils"
)

type Handler struct {
	service *Service
}

type CreateCustomerRequest struct {
	OfficeID             uint                       `json:"office_id" binding:"required"`
	Name                 string                     `json:"name" binding:"required"`
	BillingName          string                     `json:"billing_name" binding:"required"`
	MobileNo             string                     `json:"mobile_no" binding:"required"`
	Email                string                     `json:"email"`
	CreditDays           int                        `json:"credit_days"`
	MaritalStatus        string                     `json:"marital_status"`
	MarriageDate         string                     `json:"marriage_date"`
	BirthdayDate         string                     `json:"birthday_date"`
	Gender               string                     `json:"gender"`
	NID                  string                     `json:"nid"`
	PartnerGroupID      *uint                      `json:"partner_group_id"`
	PartnerSubGroupID   *uint                      `json:"partner_sub_group_id"`
	CreditLimit         float64                    `json:"credit_limit"`
	DefaultSalesRepID   *uint                      `json:"default_sales_rep_id"`
	Tolerance           float64                    `json:"tolerance"`
	BinNumber           string                     `json:"bin_number"`
	PriceTypeID         *uint                      `json:"price_type_id"`
	TCSPercentage       float64                    `json:"tcs_percentage"`
	VATRegNoCentral      string                     `json:"vat_reg_no_central"`
	TINNumber            string                     `json:"tin_number"`
	TradeLicenseNumber   string                     `json:"trade_license_number"`
	IsForeignCustomer    bool                       `json:"is_foreign_customer"`
	IsDistributor        bool                       `json:"is_distributor"`
	IsEmployeeCustomer   bool                       `json:"is_employee_customer"`
	UserID               *uint                      `json:"user_id"`
	BillingContactPerson string                     `json:"billing_contact_person"`
	BillingContactNo     string                     `json:"billing_contact_no"`
	DistrictID           *uint                      `json:"district_id"`
	ThanaID              *uint                      `json:"thana_id"`
	BillingAddress       string                     `json:"billing_address" binding:"required"`
	PointExpiryDays      int                        `json:"point_expiry_days"`
	TaxBracketID         *uint                      `json:"tax_bracket_id"`
	Remarks              string                     `json:"remarks"`
	Attachment           string                     `json:"attachment"`
	ShippingAddresses    []ShippingAddressRequest   `json:"shipping_addresses"`
	BankInfos            []BankInfoRequest          `json:"bank_infos"`
}

type ShippingAddressRequest struct {
	Code          string `json:"code"`
	RefOutlet     string `json:"ref_outlet"`
	Address       string `json:"address" binding:"required"`
	GPS           string `json:"gps"`
	ContactPerson string `json:"contact_person"`
	ContactNo     string `json:"contact_no"`
	VATRegNo      string `json:"vat_reg_no"`
	DistrictID    *uint  `json:"district_id"`
	ThanaID       *uint  `json:"thana_id"`
}

type BankInfoRequest struct {
	BankName    string `json:"bank_name"`
	AccountNo   string `json:"account_no"`
	AccountName string `json:"account_name"`
	BranchName  string `json:"branch_name"`
	RoutingNo   string `json:"routing_no"`
}

type UpdateCustomerRequest struct {
	OfficeID             uint                       `json:"office_id" binding:"required"`
	Name                 string                     `json:"name" binding:"required"`
	BillingName          string                     `json:"billing_name" binding:"required"`
	MobileNo             string                     `json:"mobile_no" binding:"required"`
	Email                string                     `json:"email"`
	CreditDays           int                        `json:"credit_days"`
	MaritalStatus        string                     `json:"marital_status"`
	MarriageDate         string                     `json:"marriage_date"`
	BirthdayDate         string                     `json:"birthday_date"`
	Gender               string                     `json:"gender"`
	NID                  string                     `json:"nid"`
	PartnerGroupID      *uint                      `json:"partner_group_id"`
	PartnerSubGroupID   *uint                      `json:"partner_sub_group_id"`
	CreditLimit         float64                    `json:"credit_limit"`
	DefaultSalesRepID   *uint                      `json:"default_sales_rep_id"`
	Tolerance           float64                    `json:"tolerance"`
	BinNumber           string                     `json:"bin_number"`
	PriceTypeID         *uint                      `json:"price_type_id"`
	TCSPercentage       float64                    `json:"tcs_percentage"`
	VATRegNoCentral      string                     `json:"vat_reg_no_central"`
	TINNumber            string                     `json:"tin_number"`
	TradeLicenseNumber   string                     `json:"trade_license_number"`
	IsForeignCustomer    bool                       `json:"is_foreign_customer"`
	IsDistributor        bool                       `json:"is_distributor"`
	IsEmployeeCustomer   bool                       `json:"is_employee_customer"`
	UserID               *uint                      `json:"user_id"`
	BillingContactPerson string                     `json:"billing_contact_person"`
	BillingContactNo     string                     `json:"billing_contact_no"`
	DistrictID           *uint                      `json:"district_id"`
	ThanaID              *uint                      `json:"thana_id"`
	BillingAddress       string                     `json:"billing_address" binding:"required"`
	PointExpiryDays      int                        `json:"point_expiry_days"`
	TaxBracketID         *uint                      `json:"tax_bracket_id"`
	Remarks              string                     `json:"remarks"`
	Attachment           string                     `json:"attachment"`
	ShippingAddresses    []ShippingAddressRequest   `json:"shipping_addresses"`
	BankInfos            []BankInfoRequest          `json:"bank_infos"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateCustomerRequest true "customer payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/customers [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	c := &Customer{
		OfficeID:             req.OfficeID,
		Name:                 req.Name,
		BillingName:          req.BillingName,
		MobileNo:             req.MobileNo,
		Email:                req.Email,
		CreditDays:           req.CreditDays,
		MaritalStatus:        req.MaritalStatus,
		Gender:               req.Gender,
		NID:                  req.NID,
		PartnerGroupID:      req.PartnerGroupID,
		PartnerSubGroupID:   req.PartnerSubGroupID,
		CreditLimit:         req.CreditLimit,
		DefaultSalesRepID:   req.DefaultSalesRepID,
		Tolerance:           req.Tolerance,
		BinNumber:           req.BinNumber,
		PriceTypeID:         req.PriceTypeID,
		TCSPercentage:       req.TCSPercentage,
		VATRegNoCentral:      req.VATRegNoCentral,
		TINNumber:            req.TINNumber,
		TradeLicenseNumber:   req.TradeLicenseNumber,
		IsForeignCustomer:    req.IsForeignCustomer,
		IsDistributor:        req.IsDistributor,
		IsEmployeeCustomer:   req.IsEmployeeCustomer,
		UserID:               req.UserID,
		BillingContactPerson: req.BillingContactPerson,
		BillingContactNo:     req.BillingContactNo,
		DistrictID:           req.DistrictID,
		ThanaID:              req.ThanaID,
		BillingAddress:       req.BillingAddress,
		PointExpiryDays:      req.PointExpiryDays,
		TaxBracketID:         req.TaxBracketID,
		Remarks:              req.Remarks,
		Attachment:           req.Attachment,
	}

	if req.MarriageDate != "" {
		mDate, err := utils.ParseDate(req.MarriageDate)
		if err == nil {
			c.MarriageDate = &mDate
		}
	}

	if req.BirthdayDate != "" {
		bDate, err := utils.ParseDate(req.BirthdayDate)
		if err == nil {
			c.BirthdayDate = &bDate
		}
	}

	for _, sa := range req.ShippingAddresses {
		c.ShippingAddresses = append(c.ShippingAddresses, CustomerShippingAddress{
			Code:          sa.Code,
			RefOutlet:     sa.RefOutlet,
			Address:       sa.Address,
			GPS:           sa.GPS,
			ContactPerson: sa.ContactPerson,
			ContactNo:     sa.ContactNo,
			VATRegNo:      sa.VATRegNo,
			DistrictID:    sa.DistrictID,
			ThanaID:       sa.ThanaID,
		})
	}

	for _, bi := range req.BankInfos {
		c.BankInfos = append(c.BankInfos, CustomerBankInfo{
			BankName:    bi.BankName,
			AccountNo:   bi.AccountNo,
			AccountName: bi.AccountName,
			BranchName:  bi.BranchName,
			RoutingNo:   bi.RoutingNo,
		})
	}

	if err := h.service.Create(ctx.Request.Context(), c); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create customer", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "customer created successfully", c)
}

// List godoc
// @Summary List all customers
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /partners/customers [get]
func (h *Handler) List(ctx *gin.Context) {
	customers, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch customers", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "customers fetched successfully", customers)
}

// Get godoc
// @Summary Get customer by ID
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/customers/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	customer, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "customer not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "customer fetched successfully", customer)
}

// Update godoc
// @Summary Update customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param payload body UpdateCustomerRequest true "customer payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/customers/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	c := &Customer{
		ID:                   uint(id),
		OfficeID:             req.OfficeID,
		Name:                 req.Name,
		BillingName:          req.BillingName,
		MobileNo:             req.MobileNo,
		Email:                req.Email,
		CreditDays:           req.CreditDays,
		MaritalStatus:        req.MaritalStatus,
		Gender:               req.Gender,
		NID:                  req.NID,
		PartnerGroupID:      req.PartnerGroupID,
		PartnerSubGroupID:   req.PartnerSubGroupID,
		CreditLimit:         req.CreditLimit,
		DefaultSalesRepID:   req.DefaultSalesRepID,
		Tolerance:           req.Tolerance,
		BinNumber:           req.BinNumber,
		PriceTypeID:         req.PriceTypeID,
		TCSPercentage:       req.TCSPercentage,
		VATRegNoCentral:      req.VATRegNoCentral,
		TINNumber:            req.TINNumber,
		TradeLicenseNumber:   req.TradeLicenseNumber,
		IsForeignCustomer:    req.IsForeignCustomer,
		IsDistributor:        req.IsDistributor,
		IsEmployeeCustomer:   req.IsEmployeeCustomer,
		UserID:               req.UserID,
		BillingContactPerson: req.BillingContactPerson,
		BillingContactNo:     req.BillingContactNo,
		DistrictID:           req.DistrictID,
		ThanaID:              req.ThanaID,
		BillingAddress:       req.BillingAddress,
		PointExpiryDays:      req.PointExpiryDays,
		TaxBracketID:         req.TaxBracketID,
		Remarks:              req.Remarks,
		Attachment:           req.Attachment,
	}

	if req.MarriageDate != "" {
		mDate, err := utils.ParseDate(req.MarriageDate)
		if err == nil {
			c.MarriageDate = &mDate
		}
	}

	if req.BirthdayDate != "" {
		bDate, err := utils.ParseDate(req.BirthdayDate)
		if err == nil {
			c.BirthdayDate = &bDate
		}
	}

	// For simple implementation, we replace all nested items. 
	// In production, we might want more granular updates.
	for _, sa := range req.ShippingAddresses {
		c.ShippingAddresses = append(c.ShippingAddresses, CustomerShippingAddress{
			Code:          sa.Code,
			RefOutlet:     sa.RefOutlet,
			Address:       sa.Address,
			GPS:           sa.GPS,
			ContactPerson: sa.ContactPerson,
			ContactNo:     sa.ContactNo,
			VATRegNo:      sa.VATRegNo,
			DistrictID:    sa.DistrictID,
			ThanaID:       sa.ThanaID,
		})
	}

	for _, bi := range req.BankInfos {
		c.BankInfos = append(c.BankInfos, CustomerBankInfo{
			BankName:    bi.BankName,
			AccountNo:   bi.AccountNo,
			AccountName: bi.AccountName,
			BranchName:  bi.BranchName,
			RoutingNo:   bi.RoutingNo,
		})
	}

	if err := h.service.Update(ctx.Request.Context(), c); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update customer", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "customer updated successfully", c)
}

// Delete godoc
// @Summary Delete customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/customers/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete customer", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "customer deleted successfully", nil)
}
