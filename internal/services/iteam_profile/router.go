package iteam_profile

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(
	v1 *gin.RouterGroup,
	db *gorm.DB,
	rdb *redis.Client,
	h *Handlers,
	authMiddleware gin.HandlerFunc,
) {
	profileGroup := v1.Group("/profile-items", authMiddleware)

	// Category routes
	categories := profileGroup.Group("/categories")
	{
		categories.POST("", middleware.Permission(db, rdb, "categories.create"), h.Category.Create)
		categories.GET("", middleware.Permission(db, rdb, "categories.view"), h.Category.List)
		categories.GET("/:id", middleware.Permission(db, rdb, "categories.view"), h.Category.Get)
		categories.PUT("/:id", middleware.Permission(db, rdb, "categories.edit"), h.Category.Update)
		categories.DELETE("/:id", middleware.Permission(db, rdb, "categories.delete"), h.Category.Delete)
	}

	// SubCategory routes
	subCategories := profileGroup.Group("/sub-categories")
	{
		subCategories.POST("", middleware.Permission(db, rdb, "sub_categories.create"), h.SubCategory.Create)
		subCategories.GET("", middleware.Permission(db, rdb, "sub_categories.view"), h.SubCategory.List)
		subCategories.GET("/:id", middleware.Permission(db, rdb, "sub_categories.view"), h.SubCategory.Get)
		subCategories.PUT("/:id", middleware.Permission(db, rdb, "sub_categories.edit"), h.SubCategory.Update)
		subCategories.DELETE("/:id", middleware.Permission(db, rdb, "sub_categories.delete"), h.SubCategory.Delete)
	}

	// ItemType routes
	itemTypes := profileGroup.Group("/item-types")
	{
		itemTypes.POST("", middleware.Permission(db, rdb, "item_types.create"), h.ItemType.Create)
		itemTypes.GET("", middleware.Permission(db, rdb, "item_types.view"), h.ItemType.List)
		itemTypes.GET("/:id", middleware.Permission(db, rdb, "item_types.view"), h.ItemType.Get)
		itemTypes.PUT("/:id", middleware.Permission(db, rdb, "item_types.edit"), h.ItemType.Update)
		itemTypes.DELETE("/:id", middleware.Permission(db, rdb, "item_types.delete"), h.ItemType.Delete)
	}

	// Tags routes
	tagRoutes := profileGroup.Group("/tags")
	{
		tagRoutes.POST("", middleware.Permission(db, rdb, "tags.create"), h.Tags.Create)
		tagRoutes.GET("", middleware.Permission(db, rdb, "tags.view"), h.Tags.List)
		tagRoutes.GET("/:id", middleware.Permission(db, rdb, "tags.view"), h.Tags.Get)
		tagRoutes.PUT("/:id", middleware.Permission(db, rdb, "tags.edit"), h.Tags.Update)
		tagRoutes.DELETE("/:id", middleware.Permission(db, rdb, "tags.delete"), h.Tags.Delete)
	}

	// -----------Minor Category routes--------------
	minorCategories := profileGroup.Group("/minor-categories")
	{
		minorCategories.POST("", middleware.Permission(db, rdb, "minor_categories.create"), h.MinorCategory.Create)
		minorCategories.GET("", middleware.Permission(db, rdb, "minor_categories.view"), h.MinorCategory.List)
		minorCategories.GET("/:id", middleware.Permission(db, rdb, "minor_categories.view"), h.MinorCategory.Get)
		minorCategories.PUT("/:id", middleware.Permission(db, rdb, "minor_categories.edit"), h.MinorCategory.Update)
		minorCategories.DELETE("/:id", middleware.Permission(db, rdb, "minor_categories.delete"), h.MinorCategory.Delete)
	}

	// Sales Setup routes
	salesSetups := profileGroup.Group("/sales-setups")
	{
		salesSetups.POST("", middleware.Permission(db, rdb, "sales_setups.create"), h.SalesSetup.Create)
		salesSetups.GET("", middleware.Permission(db, rdb, "sales_setups.view"), h.SalesSetup.List)
		salesSetups.GET("/:id", middleware.Permission(db, rdb, "sales_setups.view"), h.SalesSetup.Get)
		salesSetups.PUT("/:id", middleware.Permission(db, rdb, "sales_setups.edit"), h.SalesSetup.Update)
		salesSetups.DELETE("/:id", middleware.Permission(db, rdb, "sales_setups.delete"), h.SalesSetup.Delete)
	}

	//-------------- UOM routes ---------------
	uom := profileGroup.Group("/uom")
	{
		uom.POST("", middleware.Permission(db, rdb, "uom.create"), h.Uom.Create)
		uom.GET("", middleware.Permission(db, rdb, "uom.view"), h.Uom.List)
		uom.GET("/:id", middleware.Permission(db, rdb, "uom.view"), h.Uom.Get)
		uom.PUT("/:id", middleware.Permission(db, rdb, "uom.edit"), h.Uom.Update)
		uom.DELETE("/:id", middleware.Permission(db, rdb, "uom.delete"), h.Uom.Delete)
	}
	//------------ Sales Supply Type routes -------------
	salesSupplyTypes := profileGroup.Group("/sales-supply-types")
	{
		salesSupplyTypes.POST("", middleware.Permission(db, rdb, "sales_supply_types.create"), h.SalesSupplyType.Create)		
		salesSupplyTypes.GET("", middleware.Permission(db, rdb, "sales_supply_types.view"), h.SalesSupplyType.List)
		salesSupplyTypes.GET("/:id", middleware.Permission(db, rdb, "sales_supply_types.view"), h.SalesSupplyType.Get)
		salesSupplyTypes.PUT("/:id", middleware.Permission(db, rdb, "sales_supply_types.edit"), h.SalesSupplyType.Update)
		salesSupplyTypes.DELETE("/:id", middleware.Permission(db, rdb, "sales_supply_types.delete"), h.SalesSupplyType.Delete)
	}

	// ----------------- Sales Tax Setup routes ----------------
	salesTaxSetups := profileGroup.Group("/sales-tax-setups")
	{
		salesTaxSetups.POST("", middleware.Permission(db, rdb, "sales_tax_setups.create"), h.SalesTaxSetup.Create)		
		salesTaxSetups.GET("", middleware.Permission(db, rdb, "sales_tax_setups.view"), h.SalesTaxSetup.List)
		salesTaxSetups.GET("/:id", middleware.Permission(db, rdb, "sales_tax_setups.view"), h.SalesTaxSetup.Get)
		salesTaxSetups.PUT("/:id", middleware.Permission(db, rdb, "sales_tax_setups.edit"), h.SalesTaxSetup.Update)
		salesTaxSetups.DELETE("/:id", middleware.Permission(db, rdb, "sales_tax_setups.delete"), h.SalesTaxSetup.Delete)
	}

	//------------- Create Items routes----------------
	items := profileGroup.Group("/items")
	{
		items.POST("", middleware.Permission(db, rdb, "items.create"), h.Items.Create)
		items.GET("", middleware.Permission(db, rdb, "items.view"), h.Items.List)
		items.GET("/:id", middleware.Permission(db, rdb, "items.view"), h.Items.Get)
		items.PATCH("/:id", middleware.Permission(db, rdb, "items.edit"), h.Items.Update)
		items.DELETE("/:id", middleware.Permission(db, rdb, "items.delete"), h.Items.Delete)
	}
}
