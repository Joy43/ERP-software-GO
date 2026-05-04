package partner

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
	partnerGroup := v1.Group("/partners", authMiddleware)

	// Supplier routes
	suppliers := partnerGroup.Group("/suppliers")
	{
		suppliers.POST("", middleware.Permission(db, rdb, "suppliers.create"), h.Supplier.Create)
		suppliers.GET("", middleware.Permission(db, rdb, "suppliers.view"), h.Supplier.List)
		suppliers.GET("/:id", middleware.Permission(db, rdb, "suppliers.view"), h.Supplier.Get)
		suppliers.PUT("/:id", middleware.Permission(db, rdb, "suppliers.edit"), h.Supplier.Update)
		suppliers.DELETE("/:id", middleware.Permission(db, rdb, "suppliers.delete"), h.Supplier.Delete)
	}

	// Customer routes
	customers := partnerGroup.Group("/customers")
	{
		customers.POST("", middleware.Permission(db, rdb, "customers.create"), h.Customer.Create)
		customers.GET("", middleware.Permission(db, rdb, "customers.view"), h.Customer.List)
		customers.GET("/:id", middleware.Permission(db, rdb, "customers.view"), h.Customer.Get)
		customers.PUT("/:id", middleware.Permission(db, rdb, "customers.edit"), h.Customer.Update)
		customers.DELETE("/:id", middleware.Permission(db, rdb, "customers.delete"), h.Customer.Delete)
	}

	// Partner Group routes
	groups := partnerGroup.Group("/partner-groups")
	{
		groups.POST("", middleware.Permission(db, rdb, "partner_groups.create"), h.PartnerGroup.Create)
		groups.GET("", middleware.Permission(db, rdb, "partner_groups.view"), h.PartnerGroup.List)
		groups.GET("/:id", middleware.Permission(db, rdb, "partner_groups.view"), h.PartnerGroup.Get)
		groups.PUT("/:id", middleware.Permission(db, rdb, "partner_groups.edit"), h.PartnerGroup.Update)
		groups.DELETE("/:id", middleware.Permission(db, rdb, "partner_groups.delete"), h.PartnerGroup.Delete)
	}

	// Partner Sub Group routes
	subGroups := partnerGroup.Group("/partner-sub-groups")
	{
		subGroups.POST("", middleware.Permission(db, rdb, "partner_sub_groups.create"), h.PartnerSubGroup.Create)
		subGroups.GET("", middleware.Permission(db, rdb, "partner_sub_groups.view"), h.PartnerSubGroup.List)
		subGroups.GET("/:id", middleware.Permission(db, rdb, "partner_sub_groups.view"), h.PartnerSubGroup.Get)
		subGroups.PUT("/:id", middleware.Permission(db, rdb, "partner_sub_groups.edit"), h.PartnerSubGroup.Update)
		subGroups.DELETE("/:id", middleware.Permission(db, rdb, "partner_sub_groups.delete"), h.PartnerSubGroup.Delete)
	}

	// District routes
	districts := partnerGroup.Group("/districts")
	{
		districts.POST("", middleware.Permission(db, rdb, "districts.create"), h.District.Create)
		districts.GET("", middleware.Permission(db, rdb, "districts.view"), h.District.List)
		districts.GET("/:id", middleware.Permission(db, rdb, "districts.view"), h.District.Get)
		districts.PUT("/:id", middleware.Permission(db, rdb, "districts.edit"), h.District.Update)
		districts.DELETE("/:id", middleware.Permission(db, rdb, "districts.delete"), h.District.Delete)
	}

	// Thana routes
	thanas := partnerGroup.Group("/thanas")
	{
		thanas.POST("", middleware.Permission(db, rdb, "thanas.create"), h.Thana.Create)
		thanas.GET("", middleware.Permission(db, rdb, "thanas.view"), h.Thana.List)
		thanas.GET("/:id", middleware.Permission(db, rdb, "thanas.view"), h.Thana.Get)
		thanas.PUT("/:id", middleware.Permission(db, rdb, "thanas.edit"), h.Thana.Update)
		thanas.DELETE("/:id", middleware.Permission(db, rdb, "thanas.delete"), h.Thana.Delete)
	}

	// Tax Bracket routes
	taxBrackets := partnerGroup.Group("/tax-brackets")
	{
		taxBrackets.POST("", middleware.Permission(db, rdb, "tax_brackets.create"), h.TaxBracket.Create)
		taxBrackets.GET("", middleware.Permission(db, rdb, "tax_brackets.view"), h.TaxBracket.List)
		taxBrackets.GET("/:id", middleware.Permission(db, rdb, "tax_brackets.view"), h.TaxBracket.Get)
		taxBrackets.PUT("/:id", middleware.Permission(db, rdb, "tax_brackets.edit"), h.TaxBracket.Update)
		taxBrackets.DELETE("/:id", middleware.Permission(db, rdb, "tax_brackets.delete"), h.TaxBracket.Delete)
	}

	// Price Type routes
	priceTypes := partnerGroup.Group("/price-types")
	{
		priceTypes.POST("", middleware.Permission(db, rdb, "price_types.create"), h.PriceType.Create)
		priceTypes.GET("", middleware.Permission(db, rdb, "price_types.view"), h.PriceType.List)
		priceTypes.GET("/:id", middleware.Permission(db, rdb, "price_types.view"), h.PriceType.Get)
		priceTypes.PUT("/:id", middleware.Permission(db, rdb, "price_types.edit"), h.PriceType.Update)
		priceTypes.DELETE("/:id", middleware.Permission(db, rdb, "price_types.delete"), h.PriceType.Delete)
	}

	// Sales Representative routes
	salesReps := partnerGroup.Group("/sales-representatives")
	{
		salesReps.POST("", middleware.Permission(db, rdb, "sales_representatives.create"), h.SalesRepresentative.Create)
		salesReps.GET("", middleware.Permission(db, rdb, "sales_representatives.view"), h.SalesRepresentative.List)
		salesReps.GET("/:id", middleware.Permission(db, rdb, "sales_representatives.view"), h.SalesRepresentative.Get)
		salesReps.PUT("/:id", middleware.Permission(db, rdb, "sales_representatives.edit"), h.SalesRepresentative.Update)
		salesReps.DELETE("/:id", middleware.Permission(db, rdb, "sales_representatives.delete"), h.SalesRepresentative.Delete)
	}
}
