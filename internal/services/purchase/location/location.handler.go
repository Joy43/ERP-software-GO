package location

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	service *Service
	rdb     *redis.Client
}

func NewHandler(service *Service, rdb *redis.Client) *Handler {
	return &Handler{service: service, rdb: rdb}
}

//
// =============================
// GET ALL LOCATIONS
// =============================

// GetAllLocations godoc
// @Summary Get all locations 
// @Tags purchase--------------------------Location
// @Produce json
// @Success 200 {array} LocationWithRelationsResponse
// @Router /purchase/locations [get]
func (h *Handler) GetAllLocations(c *gin.Context) {
	data, err := h.service.GetAllLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

//
// =============================
// CREATE LOCATION
// =============================

// CreateLocation godoc
// @Summary Create a new location || create root location (without parent  "parent_id": null,)
// @Description Create location with required and optional fields
// @Tags purchase--------------------------Location
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateLocationRequest true "location creation payload"
// @Success 201 {object} LocationWithRelationsResponse
// @Router /purchase/locations [post]
func (h *Handler) CreateLocation(c *gin.Context) {
	var req CreateLocationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}


	if err := h.service.ValidateCreateLocationRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err := h.service.CreateLocation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, data)
}


//
// =============================
// UPDATE LOCATION
// =============================

// UpdateLocation godoc
// @Summary Update location
// @Description Update location by ID
// @Tags purchase--------------------------Location
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Location ID"
// @Param payload body UpdateLocationRequest true "location update payload"
// @Success 200 {object} LocationWithRelationsResponse
// @Router /purchase/locations/{id} [patch]
func (h *Handler) UpdateLocation(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid location id",
		})
		return
	}

	var req UpdateLocationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	data, err := h.service.UpdateLocation(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
// ----------------
// DELETE LOCATION (SOFT DELETE)
//-------------------	
// DeleteLocation godoc
// @Summary Delete location
// @Description Soft delete location by ID
// @Tags purchase--------------------------Location
// @Produce json
// @Security BearerAuth
// @Param id path int true "Location ID"
// @Success 200 {object} map[string]interface{}
// @Router /purchase/locations/{id} [delete]
func (h *Handler) DeleteLocation(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid location id",
		})
		return
	}

	err = h.service.DeleteLocation(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "location deleted successfully",
	})
}
