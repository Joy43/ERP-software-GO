package iteam_profile

import (
	salessupplytype "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/SalesSupplyType"
	salestaxsetup "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/SalesTaxSetup"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/item_type"
	itemsmodule "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/minor_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sales_setup"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sub_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/tags"
	umo_measurement "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/uom"
)

//-------  Handlers groups all item_profile-related handlers----------
type Handlers struct {
	Category            *category.Handler
	SubCategory         *sub_category.Handler
	Tags                *tags.Handler
	ItemType            *item_type.Handler
	MinorCategory       *minor_category.Handler
	SalesSetup          *sales_setup.Handler
	Uom                 *umo_measurement.Handler
	SalesSupplyType     *salessupplytype.Handler
	SalesTaxSetup       *salestaxsetup.Handler
	Items               *itemsmodule.Handler
	
}
