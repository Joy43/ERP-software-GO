package grn

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 32)
	return uint(val), err
}

func getCurrentUserID(c *gin.Context) *uint {
	if raw, exists := c.Get("user_id"); exists {
		switch v := raw.(type) {
		case uint:
			return &v
		case float64:
			uid := uint(v)
			return &uid
		}
	}
	return nil
}

// GetAll godoc
// @Summary      List GRNs
// @Tags         purchase--------------------- GRN
// @Produce      json
// @Security     BearerAuth
// @Param        page         query int    false "Page"      default(1)
// @Param        page_size    query int    false "Page size" default(10)
// @Param        status       query string false "Filter by status"
// @Param        receive_type query string false "DIRECT or AGAINST_PO"
// @Param        supplier_id  query uint   false "Filter by supplier"
// @Param        po_id        query uint   false "Filter by PO"
// @Param        search       query string false "Search grn_number / challan_no"
// @Success      200 {object} PaginatedGRNResponse
// @Router       /purchase/grn [get]
func (h *Handler) GetAll(c *gin.Context) {
	var f ListGRNRequest
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.GetAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetByID godoc
// @Summary      Get GRN by ID
// @Tags         purchase--------------------- GRN
// @Produce      json
// @Security     BearerAuth
// @Param        id path uint true "GRN ID"
// @Success      200 {object} GRNResponse
// @Router       /purchase/grn/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	g, err := h.service.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "GRN not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

// Create godoc
// @Summary      Create a GRN (saved as Pending, call /approve to confirm and update stock)
// @Tags         purchase--------------------- GRN
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateGRNRequest true "GRN payload"
// @Success      201 {object} GRNResponse
// @Router       /purchase/grn [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateGRNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := h.service.Create(&req, getCurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

// Approve godoc
// @Summary      Approve a PENDING GRN (confirms receipt and updates location_stocks)
// @Tags         purchase--------------------- GRN
// @Produce      json
// @Security     BearerAuth
// @Param        id path uint true "GRN ID"
// @Success      200 {object} GRNResponse
// @Router       /purchase/grn/{id}/approve [patch]
func (h *Handler) Approve(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	g, err := h.service.Approve(id, getCurrentUserID(c))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "GRN not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}
