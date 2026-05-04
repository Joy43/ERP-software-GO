package purchase

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(
	V1 *gin.RouterGroup,
	db *gorm.DB,
	rdb *redis.Client,
	h *Handlers,
	authMiddleware gin.HandlerFunc,
) {
	purchaseGroup := V1.Group("/purchase", authMiddleware)

	// ============================
	// LOCATION ROUTES
	// ============================
	locationRoutes := purchaseGroup.Group("/locations")

	{
		locationRoutes.GET("", middleware.Permission(db, rdb, "locations.view"), h.Location.GetAllLocations)
		locationRoutes.POST("", middleware.Permission(db, rdb, "locations.create"), h.Location.CreateLocation)
		locationRoutes.PATCH("/:id", middleware.Permission(db, rdb, "locations.update"), h.Location.UpdateLocation)
		locationRoutes.DELETE("/:id", middleware.Permission(db, rdb, "locations.delete"), h.Location.DeleteLocation)
		
	}


	// ============================
	//----------  Inventory Types ROUTES END -----------
	// ============================
	InventoryRoutes := purchaseGroup.Group("/inventory-types")

	{
		InventoryRoutes.GET("", middleware.Permission(db, rdb, "inventory.view"), h.InventoryTypes.GetAll)
		InventoryRoutes.POST("", middleware.Permission(db, rdb, "inventory.create"), h.InventoryTypes.Create)
		InventoryRoutes.PATCH("/:id", middleware.Permission(db, rdb, "inventory.update"), h.InventoryTypes.Update)
		InventoryRoutes.DELETE("/:id", middleware.Permission(db, rdb, "inventory.delete"), h.InventoryTypes.Delete)

	}

	// ============================
	// PROJECTS ROUTES
	// ============================
	projectRoutes := purchaseGroup.Group("/projects")
	{
		projectRoutes.GET("", middleware.Permission(db, rdb, "projects.view"), h.Projects.ListProjects)
		projectRoutes.POST("", middleware.Permission(db, rdb, "projects.create"), h.Projects.CreateProject)
		projectRoutes.PUT("/:id", middleware.Permission(db, rdb, "projects.update"), h.Projects.UpdateProject)
		projectRoutes.PATCH("/:id", middleware.Permission(db, rdb, "projects.update"), h.Projects.UpdateProject)
		projectRoutes.DELETE("/:id", middleware.Permission(db, rdb, "projects.delete"), h.Projects.DeleteProject)
		projectRoutes.PUT("/:id/status", middleware.Permission(db, rdb, "projects.update"), h.Projects.ChangeProjectStatus)
	}


	// ============================
	// PURCHASE ORDER ROUTES
	// ============================
	poRoutes := purchaseGroup.Group("/orders")
	{
		poRoutes.GET("", middleware.Permission(db, rdb, "purchase_orders.view"), h.PurchaseOrder.GetAll)
		poRoutes.GET("/ordered-requisitions", middleware.Permission(db, rdb, "purchase_orders.view"), h.PurchaseOrder.GetOrderedRequisitions)
		poRoutes.POST("", middleware.Permission(db, rdb, "purchase_orders.create"), h.PurchaseOrder.Create)
		poRoutes.GET("/:id", middleware.Permission(db, rdb, "purchase_orders.view"), h.PurchaseOrder.GetByID)
		poRoutes.PATCH("/:id", middleware.Permission(db, rdb, "purchase_orders.update"), h.PurchaseOrder.Update)
		poRoutes.PATCH("/:id/status", middleware.Permission(db, rdb, "purchase_orders.approve"), h.PurchaseOrder.UpdateStatus)
		poRoutes.DELETE("/:id", middleware.Permission(db, rdb, "purchase_orders.delete"), h.PurchaseOrder.Delete)
	}

	// ============================
	// GRN ROUTES
	// ============================
	grnRoutes := purchaseGroup.Group("/grn")
	{
		grnRoutes.GET("", middleware.Permission(db, rdb, "grn.view"), h.GRN.GetAll)
		grnRoutes.POST("", middleware.Permission(db, rdb, "grn.create"), h.GRN.Create)
		grnRoutes.GET("/:id", middleware.Permission(db, rdb, "grn.view"), h.GRN.GetByID)
		grnRoutes.PATCH("/:id/approve", middleware.Permission(db, rdb, "grn.approve"), h.GRN.Approve)
	}

	// ============================
	// REQUISITION ROUTES
	// ============================
	requisitionRoutes := purchaseGroup.Group("/requisitions")
	{
		//------- List & Create------------
		requisitionRoutes.GET("", middleware.Permission(db, rdb, "requisitions.view"), h.Requisitions.GetAllRequisitions)
		requisitionRoutes.POST("", middleware.Permission(db, rdb, "requisitions.create"), h.Requisitions.CreateRequisition)
 
		// --------- Single record operations ---------
		requisitionRoutes.GET("/:id", middleware.Permission(db, rdb, "requisitions.view"), h.Requisitions.GetRequisitionByID)
		requisitionRoutes.PATCH("/:id", middleware.Permission(db, rdb, "requisitions.update"), h.Requisitions.UpdateRequisition)
		requisitionRoutes.DELETE("/:id", middleware.Permission(db, rdb, "requisitions.delete"), h.Requisitions.DeleteRequisition)

		// ------------ Status transition -----------
		requisitionRoutes.PATCH("/:id/status", middleware.Permission(db, rdb, "requisitions.approve"), h.Requisitions.UpdateRequisitionStatus)
		requisitionRoutes.GET("/:id/history", middleware.Permission(db, rdb, "requisitions.view"), h.Requisitions.GetRequisitionHistory)
		requisitionRoutes.GET("/summary", middleware.Permission(db, rdb, "requisitions.view"), h.Requisitions.GetRequisitionSummary)
	}
// ============================
	// PURCHASE PAYMENTS ROUTES
	// ============================
	paymentByGRNRoutes := purchaseGroup.Group("/payment-by-grn")
	{
		paymentByGRNRoutes.POST("", middleware.Permission(db, rdb, "purchase_payments.create"), h.PurchasePayments.CreatePaymentByGRN)
		paymentByGRNRoutes.GET("", middleware.Permission(db, rdb, "purchase_payments.view"), h.PurchasePayments.GetAllPaymentByGRN)
	}
 
	advancePaymentRoutes := purchaseGroup.Group("/advance-payments")
	{
		advancePaymentRoutes.POST("", middleware.Permission(db, rdb, "purchase_payments.create"), h.PurchasePayments.CreateAdvancePayment)
		advancePaymentRoutes.GET("", middleware.Permission(db, rdb, "purchase_payments.view"), h.PurchasePayments.GetAllAdvancePayments)
	}
 
	supplierBillRoutes := purchaseGroup.Group("/supplier-bills")
	{
		supplierBillRoutes.POST("", middleware.Permission(db, rdb, "purchase_payments.create"), h.PurchasePayments.CreateSupplierBill)
		supplierBillRoutes.GET("", middleware.Permission(db, rdb, "purchase_payments.view"), h.PurchasePayments.GetAllSupplierBills)
	}

	// ========== 
	// PURCHASE RECEIVE ROUTES
	// ==========
	// purchaseReceiveRoutes := purchaseGroup.Group("/purchase-receive")
	// {
	// 	purchaseReceiveRoutes.POST("", middleware.Permission(db, rdb, "purchase_receive.create"), h.PurchaseReceive.Create)
	// 	purchaseReceiveRoutes.GET("", middleware.Permission(db, rdb, "purchase_receive.view"), h.PurchaseReceive.GetAll)
	// 	purchaseReceiveRoutes.GET("/:id", middleware.Permission(db, rdb, "purchase_receive.view"), h.PurchaseReceive.GetByID)
	// 	purchaseReceiveRoutes.PUT("/:id", middleware.Permission(db, rdb, "purchase_receive.update"), h.PurchaseReceive.Update)
	// 	purchaseReceiveRoutes.PATCH("/:id", middleware.Permission(db, rdb, "purchase_receive.update"), h.PurchaseReceive.Update)
	// 	purchaseReceiveRoutes.DELETE("/:id", middleware.Permission(db, rdb, "purchase_receive.delete"), h.PurchaseReceive.Delete)
	// }

	// ========= 
	// PURCHASE RETURN ROUTES
	// ==========
	purchaseReturnRoutes := purchaseGroup.Group("/purchase-returns")
	{
		purchaseReturnRoutes.POST("", middleware.Permission(db, rdb, "purchase_return.create"), h.PurchaseReturns.Create)
		purchaseReturnRoutes.GET("", middleware.Permission(db, rdb, "purchase_return.view"), h.PurchaseReturns.GetAll)

	}
}