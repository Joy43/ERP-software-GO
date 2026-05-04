package purchasereturn


import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
	rdb     *redis.Client
}

func NewHandler(service *Service, rdb *redis.Client) *Handler {
	return &Handler{service: service, rdb: rdb}
}

// GetAll godoc
// @Summary      List purchase returns (paginated + filtered)
// @Tags         purchase--------------------- Purchase Return
// @Produce      json
// @Param        page        query int    false "Page number"         default(1)
// @Param        page_size   query int    false "Page size (max 100)" default(10)
// @Param        supplier_id query uint   false "Filter by supplier"
// @Param        office_id   query uint   false "Filter by office"
// @Param        location_id query uint   false "Filter by location"
// @Param        item_id     query uint   false "Filter by item"
// @Param        search      query string false "Search in remarks / return number"
// @Param        sort_order  query string false "asc or desc"         default(desc)
// @Success      200 {object} PaginatedPurchaseReturnResponse
// @Router       /purchase/purchase-returns [get]
func (h *Handler) GetAll(c *gin.Context) {
	var filter ListPurchaseReturnRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --------Page size guard----------
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	data, err := h.service.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to fetch purchase returns",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// Create godoc
// @Summary      Create a new purchase return
// @Tags         purchase--------------------- Purchase Return
// @Accept       json
// @Produce      json
// @Param        request body CreatePurchaseReturnRequest true "Purchase return payload"
// @Success      201 {object} PurchaseReturnResponse
// @Router       /purchase/purchase-returns [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreatePurchaseReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := h.service.Create(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == gorm.ErrRecordNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error":   err.Error(),
			"message": "Failed to create purchase return",
		})
		return
	}

	c.JSON(http.StatusCreated, record)
}