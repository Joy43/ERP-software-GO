package purchaseorder

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

// GetOrderedRequisitions godoc
// @Summary      List ORDERED requisitions with their linked PO ID
// @Tags         purchase--------------------- Purchase Order
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Router       /purchase/orders/ordered-requisitions [get]
func (h *Handler) GetOrderedRequisitions(c *gin.Context) {
	rows, err := h.service.GetOrderedRequisitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": len(rows),
		"data":  rows,
	})
}

// GetAll godoc
// @Summary      List purchase orders
// @Tags         purchase--------------------- Purchase Order
// @Produce      json
// @Security     BearerAuth
// @Param        page        query int    false "Page"       default(1)
// @Param        page_size   query int    false "Page size"  default(10)
// @Param        status      query string false "Filter by status"
// @Param        supplier_id query uint   false "Filter by supplier"
// @Param        order_type  query string false "Filter by order type"
// @Param        search      query string false "Search in po_number / remarks"
// @Success      200 {object} PaginatedPOResponse
// @Router       /purchase/orders [get]
func (h *Handler) GetAll(c *gin.Context) {
	var f ListPORequest
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
// @Summary      Get purchase order by ID
// @Tags         purchase--------------------- Purchase Order
// @Produce      json
// @Security     BearerAuth
// @Param        id path uint true "PO ID"
// @Success      200 {object} POResponse
// @Router       /purchase/orders/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	po, err := h.service.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "purchase order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, po)
}

// Create godoc
// @Summary      Create a purchase order  then type selected REQUISITION_BASED("po_date": "2026-05-02","order_type": "REQUISITION_BASED","requisition_id": 4,"office_id": 1,"location_id": 9,"supplier_id": 4)
// @Tags         purchase--------------------- Purchase Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreatePORequest true "PO payload"
// @Success      201 {object} POResponse
// @Router       /purchase/orders [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreatePORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	po, err := h.service.Create(&req, getCurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, po)
}

// Update godoc
// @Summary      Update a purchase order (Pending only)
// @Tags         purchase--------------------- Purchase Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path uint          true "PO ID"
// @Param        request body UpdatePORequest true "Fields to update"
// @Success      200 {object} POResponse
// @Router       /purchase/orders/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdatePORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	po, err := h.service.Update(id, &req, getCurrentUserID(c))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "purchase order not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, po)
}

// UpdateStatus godoc
// @Summary      Transition purchase order status
// @Tags         purchase--------------------- Purchase Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path uint                true "PO ID"
// @Param        request body UpdatePOStatusRequest true "Status payload"
// @Success      200 {object} POResponse
// @Router       /purchase/orders/{id}/status [patch]
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdatePOStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	po, err := h.service.UpdateStatus(id, &req, getCurrentUserID(c))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "purchase order not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, po)
}

// Delete godoc
// @Summary      Delete a purchase order (PENDING only)
// @Tags         purchase--------------------- Purchase Order
// @Produce      json
// @Security     BearerAuth
// @Param        id path uint true "PO ID"
// @Success      200 {object} map[string]interface{}
// @Router       /purchase/orders/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "purchase order not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "purchase order deleted", "id": id})
}
