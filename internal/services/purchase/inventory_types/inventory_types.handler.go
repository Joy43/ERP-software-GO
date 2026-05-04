package inventorytypes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

//------------------------------
// 	------- CREATE --------
//------------------------------
//CreateInventoryType godoc
// @Summary Create a new inventory type
// @Description Create inventory type with required and optional fields
// @Tags purchase--------------------------InventoryType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateInventoryTypeRequest true "inventory type creation payload"
// @Success 201 {object} InventoryTypeResponse
// @Router /purchase/inventory-types [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateInventoryTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created successfully"})
}

//
// =============================
// GET ALL INVENTORY TYPES
// =============================

// GetAllInventoryTypes godoc
// @Summary Get all inventory types
// @Tags purchase--------------------------InventoryType
// @Produce json
// @Success 200 {array} InventoryTypeResponse
// @Router /purchase/inventory-types [get]
func (h *Handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}



// UPDATE
// UpdateInventoryType godoc
// @Summary Update an inventory type
// @Description Update inventory type fields (type_code, type_name, description, is_active)
// @Tags purchase--------------------------InventoryType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory Type ID"
// @Param payload body UpdateInventoryTypeRequest true "inventory type update payload"
// @Success 200 {object} InventoryTypeResponse
// @Router /purchase/inventory-types/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req UpdateInventoryTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// ------------------
// DELETE
//-------------------

// DeleteInventoryType godoc
// @Summary Delete an inventory type
// @Description Soft delete an inventory type by ID
// @Tags purchase--------------------------InventoryType
// @Security BearerAuth
// @Param id path int true "Inventory Type ID"
// @Success 200 {object} InventoryTypeResponse
// @Router /purchase/inventory-types/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}