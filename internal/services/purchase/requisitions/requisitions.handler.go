package requisitions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Handler represents the requisition handler
type Handler struct {
	service *Service
	rdb     *redis.Client
}

// NewHandler creates a new requisition handler
func NewHandler(service *Service, rdb *redis.Client) *Handler {
	return &Handler{service: service, rdb: rdb}
}

// getCurrentUserID extracts the logged-in user ID from context (set by auth middleware)
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

// ========================================
// GET ALL
// ========================================

// GetAllRequisitions godoc
// @Summary      List requisitions (paginated + filtered)
// @Tags         purchase--------------------- Requisition
// @Produce      json
// @Security     BearerAuth
// @Param        page             query int    false "Page number"          default(1)
// @Param        page_size        query int    false "Page size (max 100)"  default(10)
// @Param        status           query string false "Filter by status"
// @Param        requisition_type query string false "Filter by type"
// @Param        department_id    query uint   false "Filter by department"
// @Param        project_id       query uint   false "Filter by project"
// @Param        employee_id      query uint   false "Filter by employee"
// @Param        search           query string false "Search in remarks / description / number"
// @Param        sort_by          query string false "Field to sort by"     default(created_at)
// @Param        sort_order       query string false "asc or desc"          default(desc)
// @Success      200 {object} PaginatedRequisitionResponse
// @Router       /purchase/requisitions [get]
func (h *Handler) GetAllRequisitions(c *gin.Context) {
	var filter ListRequisitionRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.GetAllRequisitions(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to fetch requisitions",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ========================================
// GET BY ID
// ========================================

// GetRequisitionByID godoc
// @Summary      Get a single requisition
// @Tags         purchase--------------------- Requisition
// @Produce      json
// @Security     BearerAuth
// @Param        id path uint true "Requisition ID"
// @Success      200 {object} RequisitionResponse
// @Router       /purchase/requisitions/{id} [get]
func (h *Handler) GetRequisitionByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requisition ID"})
		return
	}

	requisition, err := h.service.GetRequisitionByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Requisition not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to fetch requisition",
		})
		return
	}

	c.JSON(http.StatusOK, requisition)
}

// ========================================
//---------  CREATE -----------
// ========================================

// CreateRequisition godoc
// @Summary      Create a new requisition
// @Tags         purchase--------------------- Requisition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateRequisitionRequest true "Requisition payload"
// @Success      201 {object} RequisitionResponse
// @Router       /purchase/requisitions [post]
func (h *Handler) CreateRequisition(c *gin.Context) {
	var req CreateRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getCurrentUserID(c)

	requisition, err := h.service.CreateRequisition(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create requisition",
		})
		return
	}

	c.JSON(http.StatusCreated, requisition)
}

// ========================================
//----------	 UPDATE (Pending only) ---------
// ========================================

// UpdateRequisition godoc
// @Summary      Update a requisition (Pending status only-- when item update then must be item_id is required)
// @Tags         purchase--------------------- Requisition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path uint                  true "Requisition ID"
// @Param        request body UpdateRequisitionRequest true "Fields to update"
// @Success      200 {object} RequisitionResponse
// @Router       /purchase/requisitions/{id} [patch]
func (h *Handler) UpdateRequisition(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requisition ID"})
		return
	}

	var req UpdateRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getCurrentUserID(c)

	requisition, err := h.service.UpdateRequisition(id, &req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Requisition not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to update requisition",
		})
		return
	}

	c.JSON(http.StatusOK, requisition)
}

// ========================================
// DELETE (DRAFT / CANCELLED only)
// ========================================

// DeleteRequisition godoc
// @Summary      Delete a requisition (pending or CANCELLED only)
// @Tags         purchase--------------------- Requisition
// @Produce      json
// @Security     BearerAuth
// @Param        id path uint true "Requisition ID"
// @Success      200 {object} map[string]interface{}
// @Router       /purchase/requisitions/{id} [delete]
func (h *Handler) DeleteRequisition(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requisition ID"})
		return
	}

	if err := h.service.DeleteRequisition(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Requisition not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete requisition",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Requisition deleted successfully",
		"id":      id,
	})
}

// ========================================
// UPDATE STATUS
// ========================================

// UpdateRequisitionStatus godoc
// @Summary      Transition requisition status
// @Description  Validates allowed transitions and records status history.
// @Tags         purchase--------------------- Requisition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path uint                        true "Requisition ID"
// @Param        request body UpdateRequisitionStatusRequest true "Status transition payload"
// @Success      200 {object} RequisitionResponse
// @Router       /purchase/requisitions/{id}/status [patch]
func (h *Handler) UpdateRequisitionStatus(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requisition ID"})
		return
	}

	var req UpdateRequisitionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getCurrentUserID(c)

	requisition, err := h.service.UpdateRequisitionStatus(id, &req, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Requisition not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to update requisition status",
		})
		return
	}

	c.JSON(http.StatusOK, requisition)
}

// ========================================
//------- HELPER--------
// ========================================

func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// ========================================
// GET HISTORY by id
// ========================================
// GetRequisitionHistory godoc
// @Summary      Get requisition status history
// @Tags         purchase--------------------- Requisition
// @Produce      json
// @Security     BearerAuth
// @Param        id path uint true "Requisition ID"
// @Success      200 {object} map[string]interface{}
// @Router       /purchase/requisitions/{id}/history [get]	
func (h *Handler) GetRequisitionHistory(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requisition ID"})
		return
	}

	history, err := h.service.GetRequisitionHistory(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Requisition not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requisition_id": id,
		"total":          len(history),
		"history":        history,
	})
}


// ========================================
// REQUISITION SUMMARY
// ========================================

// GetRequisitionSummary godoc
// @Summary      Get requisition summary (counts by status)
// @Tags         purchase--------------------- Requisition
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Router       /purchase/requisitions/summary [get]
func (h *Handler) GetRequisitionSummary(c *gin.Context) {
	summary, err := h.service.GetRequisitionSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
