package purchase

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/grn"
	inventorytypes "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/inventory_types"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
	purchaseorder "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/order"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/projects"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/purchasepayments"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/purchasereturn"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/requisitions"
)

type Handlers struct {
	Location       *location.Handler
	InventoryTypes *inventorytypes.Handler
	Projects       *projects.Handler
	Requisitions   *requisitions.Handler
	PurchaseOrder  *purchaseorder.Handler
	GRN            *grn.Handler
	PurchasePayments *purchasepayments.Handler
	
	PurchaseReturns   *purchasereturn.Handler
}
