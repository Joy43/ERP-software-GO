package app
import (
	"fmt"
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/config"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/database"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/router"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/audit"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/audit/transaction_history"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/department"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/designation"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/payment_mode"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/permission"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/responsibility_transfer"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/role"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/wallet"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile"
	salessupplytype "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/SalesSupplyType"
	salestaxsetup "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/SalesTaxSetup"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/item_type"
	items "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/minor_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sales_setup"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sub_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/tags"
	umo_measurement "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/uom"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/customer"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/district"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_sub_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/price_type"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/sales_representative"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/tax_bracket"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/thana"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/grn"
	inventorytypes "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/inventory_types"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
	purchaseorder "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/order"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/projects"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/purchasepayments"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/purchasereturn"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/requisitions"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/stock"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/redis"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.NewMysql(cfg.MysqlDSN)
	if err != nil {
		return err
	}

	//--------- Redis client------------
	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("setup redis: %w", err)
	}
	defer redisClient.Close()

	//----------  Permission module -----------
	permissionRepository := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepository)
	permissionHandler := permission.NewHandler(permissionService)

	// --------------- Role module ---------------
	roleRepository := role.NewRepository(db)
	roleService := role.NewService(roleRepository, permissionRepository)
	roleHandler := role.NewHandler(roleService)

	//---------------- User module ---------------
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, roleRepository, cfg)
	userHandler := user.NewHandler(userService, redisClient)

	// --------------- Designation module ---------------
	designationRepo := designation.NewRepository(db)
	designationService := designation.NewService(designationRepo)
	designationHandler := designation.NewHandler(designationService)

	// --------------- Department module ---------------
	departmentRepo := department.NewRepository(db)
	departmentService := department.NewService(departmentRepo)
	departmentHandler := department.NewHandler(departmentService)

	// --------------- Office module ---------------
	officeRepo := office.NewRepository(db)
	officeService := office.NewService(officeRepo)
	officeHandler := office.NewHandler(officeService)

	// --------------- Payment Mode module ---------------
	paymentModeRepo := payment_mode.NewRepository(db)
	paymentModeService := payment_mode.NewService(paymentModeRepo)
	paymentModeHandler := payment_mode.NewHandler(paymentModeService)

	// --------------- Responsibility Transfer module ---------------
	rtRepo := responsibility_transfer.NewRepository(db)
	rtService := responsibility_transfer.NewService(rtRepo)
	rtHandler := responsibility_transfer.NewHandler(rtService)

	// --------------- Transaction History module ---------------
	thRepo := transaction_history.NewRepository(db)
	thService := transaction_history.NewService(thRepo)
	thHandler := transaction_history.NewHandler(thService)

	// --------------- Wallet module ---------------
	walletRepo := wallet.NewRepository(db)
	walletService := wallet.NewService(walletRepo)
	walletHandler := wallet.NewHandler(walletService)

	// --------------- District module ---------------
	districtRepo := district.NewRepository(db)
	districtService := district.NewService(districtRepo)
	districtHandler := district.NewHandler(districtService)

	// --------------- Thana module ---------------
	thanaRepo := thana.NewRepository(db)
	thanaService := thana.NewService(thanaRepo)
	thanaHandler := thana.NewHandler(thanaService)

	// --------------- Partner Group module ---------------
	pgRepo := partner_group.NewRepository(db)
	pgService := partner_group.NewService(pgRepo)
	pgHandler := partner_group.NewHandler(pgService)

	// --------------- Partner Sub Group module ---------------
	psgRepo := partner_sub_group.NewRepository(db)
	psgService := partner_sub_group.NewService(psgRepo)
	psgHandler := partner_sub_group.NewHandler(psgService)

	//--------- Tax Bracket module-------------
	tbRepo := tax_bracket.NewRepository(db)
	tbService := tax_bracket.NewService(tbRepo)
	tbHandler := tax_bracket.NewHandler(tbService)

	//----------- Supplier module ---------------
	supplierRepo := supplier.NewRepository(db)
	supplierService := supplier.NewService(supplierRepo)
	supplierHandler := supplier.NewHandler(supplierService)

	// --------------- Customer module ---------------
	customerRepo := customer.NewRepository(db)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService)

	// Price Type module
	priceTypeRepo := price_type.NewRepository(db)
	priceTypeService := price_type.NewService(priceTypeRepo)
	priceTypeHandler := price_type.NewHandler(priceTypeService)

	// Sales Representative module
	salesRepRepo := sales_representative.NewRepository(db)
	salesRepService := sales_representative.NewService(salesRepRepo)
	salesRepHandler := sales_representative.NewHandler(salesRepService)

	// Category module
	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	// Sub Category module
	subCategoryRepo := sub_category.NewRepository(db)
	subCategoryService := sub_category.NewService(subCategoryRepo)
	subCategoryHandler := sub_category.NewHandler(subCategoryService)

	// Tags module
	tagsRepo := tags.NewRepository(db)
	tagsService := tags.NewService(tagsRepo)
	tagsHandler := tags.NewHandler(tagsService)

	// Minor Category module
	minorCategoryRepo := minor_category.NewRepository(db)
	minorCategoryService := minor_category.NewService(minorCategoryRepo)
	minorCategoryHandler := minor_category.NewHandler(minorCategoryService)

	// Sales Setup module
	salesSetupRepo := sales_setup.NewRepository(db)
	salesSetupService := sales_setup.NewService(salesSetupRepo)
	salesSetupHandler := sales_setup.NewHandler(salesSetupService)

	// Item Type module
	itemTypeRepo := item_type.NewRepository(db)
	itemTypeService := item_type.NewService(itemTypeRepo)
	itemTypeHandler := item_type.NewHandler(itemTypeService)

	// UOM module
	uomRepo := umo_measurement.NewRepository(db)
	uomService := umo_measurement.NewService(uomRepo)
	uomHandler := umo_measurement.NewHandler(uomService)

	// Sales Supply Type module
	salesSupplyTypeRepo := salessupplytype.NewRepository(db)
	salesSupplyTypeService := salessupplytype.NewService(salesSupplyTypeRepo)
	salesSupplyTypeHandler := salessupplytype.NewHandler(salesSupplyTypeService)

	// Sales Tax Setup module
	salesTaxSetupRepo := salestaxsetup.NewRepository(db)
	salesTaxSetupService := salestaxsetup.NewService(salesTaxSetupRepo)
	salesTaxSetupHandler := salestaxsetup.NewHandler(salesTaxSetupService)

	// Items module
	ItemRepo := items.NewRepository(db)
	ItemService := items.NewService(ItemRepo, db)
	ItemHandler := items.NewHandler(ItemService, redisClient)

	// Uploads module
	uploadsRepository := uploads.NewRepository(db)
	uploadsService := uploads.NewService(uploadsRepository)
	uploadsHandler := uploads.NewHandler(uploadsService, cfg)

	// Location module
	locationRepo := location.NewRepository(db)
	locationService := location.NewService(locationRepo)
	locationHandler := location.NewHandler(locationService, redisClient)

	// Inventory Types module
	inventoryTypeRepo := inventorytypes.NewRepository(db)
	inventoryTypeService := inventorytypes.NewService(inventoryTypeRepo)
	inventoryTypeHandler := inventorytypes.NewHandler(inventoryTypeService)

	// Projects module
	projectsRepo := projects.NewRepository(db)
	projectsService := projects.NewService(projectsRepo)
	projectsHandler := projects.NewHandler(projectsService, redisClient)

	// Purchase Order module
	poRepo := purchaseorder.NewRepository(db)

	// POCreator callback for requisitions → auto-create PO on ORDERED status
	poCreator := func(req *requisitions.Requisition, createdByID *uint) error {
		poNumber, err := poRepo.GeneratePONumber()
		if err != nil {
			return err
		}
		poItems := make([]purchaseorder.PurchaseOrderItem, 0, len(req.Items))
		var subtotal float64
		for _, ri := range req.Items {
			qty := ri.RequestQuantity
			if ri.ApprovedQuantity != nil && *ri.ApprovedQuantity > 0 {
				qty = *ri.ApprovedQuantity
			}
			unitPrice := float64(0)
			if ri.LastCost != nil {
				unitPrice = *ri.LastCost
			} else if ri.AverageCost != nil {
				unitPrice = *ri.AverageCost
			}
			lineTotal := qty * unitPrice
			subtotal += lineTotal
			riID := ri.ID
			poItems = append(poItems, purchaseorder.PurchaseOrderItem{
				ItemID:            ri.ItemID,
				RequisitionItemID: &riID,
				OrderQuantity:     qty,
				UnitPrice:         unitPrice,
				TotalAmount:       lineTotal,
			})
		}
		reqID := req.ID
		var officeID, locationID, supplierID uint
		if req.OfficeID != nil {
			officeID = *req.OfficeID
		}
		if req.LocationID != nil {
			locationID = *req.LocationID
		}
		if req.SupplierID != nil {
			supplierID = *req.SupplierID
		}
		po := &purchaseorder.PurchaseOrder{
			PONumber:      poNumber,
			PODate:        time.Now(),
			RequisitionID: &reqID,
			OrderType:     purchaseorder.POTypeRequisitionBased,
			OfficeID:      officeID,
			LocationID:    locationID,
			SupplierID:    supplierID,
			Status:        purchaseorder.POStatusPending,
			CreatedByID:   createdByID,
			Subtotal:      subtotal,
			TotalAmount:   subtotal,
			Items:         poItems,
		}
		_, err = poRepo.Create(po)
		return err
	}

	//------------- Requisitions module-----------------
	requisitionsRepo := requisitions.NewRepository(db)
	requisitionsService := requisitions.NewService(requisitionsRepo, poCreator)
	requisitionsHandler := requisitions.NewHandler(requisitionsService, redisClient)

	//------------- Stock module ---------------
	stockRepo := stock.NewRepository(db)

	// ------------- GRN module ---------------
	grnRepo := grn.NewRepository(db)
	grnService := grn.NewService(grnRepo, poRepo, stockRepo)
	grnHandler := grn.NewHandler(grnService)

	// GRNCreator callback for PO → auto-create GRN on ISSUED status
	grnCreator := func(po *purchaseorder.PurchaseOrder, createdByID *uint) error {
		grnNumber, err := grnRepo.GenerateGRNNumber()
		if err != nil {
			return err
		}
		grnItems := make([]grn.GRNItem, 0, len(po.Items))
		for _, poi := range po.Items {
			poiID := poi.ID
			lineTotal := poi.OrderQuantity * poi.UnitPrice
			vatAmt := lineTotal * poi.VatPercentage / 100
			discAmt := lineTotal * poi.DiscountPercentage / 100
			grnItems = append(grnItems, grn.GRNItem{
				ItemID:             poi.ItemID,
				POItemID:           &poiID,
				ReceivedQuantity:   poi.OrderQuantity,
				UOMID:              poi.UOMID,
				PurchasePrice:      poi.UnitPrice,
				VatPercentage:      poi.VatPercentage,
				VatAmount:          vatAmt,
				DiscountPercentage: poi.DiscountPercentage,
				DiscountAmount:     discAmt,
				TotalAmount:        lineTotal + vatAmt - discAmt,
			})
		}
		poID := po.ID
		var reqID *uint
		if po.RequisitionID != nil {
			reqID = po.RequisitionID
		}
		g := &grn.GoodsReceiptNote{
			GRNNumber:     grnNumber,
			GRNDate:       time.Now(),
			ReceiveType:   grn.GRNAgainstPO,
			Status:        grn.GRNStatusPending,
			POID:          &poID,
			RequisitionID: reqID,
			OfficeID:      po.OfficeID,
			LocationID:    po.LocationID,
			SupplierID:    po.SupplierID,
			CreatedByID:   createdByID,
			Items:         grnItems,
		}
		_, err = grnRepo.Create(g)
		return err
	}

	// Purchase Order service with GRN creator
	poService := purchaseorder.NewService(poRepo, requisitionsRepo, grnCreator)
	poHandler := purchaseorder.NewHandler(poService)
	

	// Purchase Payment module
	purchasepaymentsRepo := purchasepayments.NewRepository(db)
	purchasepaymentsService := purchasepayments.NewService(purchasepaymentsRepo)
	purchasepaymentsHandler := purchasepayments.NewHandler(purchasepaymentsService)

	// Purchase Return module
	purchaseReturnRepo := purchasereturn.NewRepository(db)
	purchaseReturnService := purchasereturn.NewService(purchaseReturnRepo)
	purchaseReturnHandler := purchasereturn.NewHandler(purchaseReturnService, redisClient)

	// Build handlers structs
	allHandlers := &router.Handlers{
		Auth: &auth.Handlers{
			User:                   userHandler,
			Role:                   roleHandler,
			Permission:             permissionHandler,
			Designation:            designationHandler,
			Department:             departmentHandler,
			Office:                 officeHandler,
			PaymentMode:            paymentModeHandler,
			ResponsibilityTransfer: rtHandler,
			Wallet:                 walletHandler,
		},
		Audit: &audit.Handlers{
			TransactionHistory: thHandler,
		},
		Partner: &partner.Handlers{
			Supplier:            supplierHandler,
			PartnerGroup:        pgHandler,
			PartnerSubGroup:     psgHandler,
			District:            districtHandler,
			Thana:               thanaHandler,
			TaxBracket:          tbHandler,
			Customer:            customerHandler,
			PriceType:           priceTypeHandler,
			SalesRepresentative: salesRepHandler,
		},
		ITeamProfile: &iteam_profile.Handlers{
			Category:        categoryHandler,
			SubCategory:     subCategoryHandler,
			Tags:            tagsHandler,
			ItemType:        itemTypeHandler,
			MinorCategory:   minorCategoryHandler,
			SalesSetup:      salesSetupHandler,
			Uom:             uomHandler,
			SalesSupplyType: salesSupplyTypeHandler,
			SalesTaxSetup:   salesTaxSetupHandler,
			Items:           ItemHandler,
		},
		Uploads: uploadsHandler,
		Purchase: &purchase.Handlers{
			Location:       locationHandler,
			InventoryTypes: inventoryTypeHandler,
			Projects:       projectsHandler,
			Requisitions:   requisitionsHandler,
			PurchaseOrder:  poHandler,
			GRN:            grnHandler,
			PurchasePayments: purchasepaymentsHandler,
			
			PurchaseReturns:  purchaseReturnHandler,
		},
	}

	r := router.New(db, redisClient, cfg, allHandlers)
	addr := fmt.Sprintf(":%s", cfg.Port)

	return r.Run(addr)
}
