package partner

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/customer"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/district"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_sub_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/tax_bracket"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/thana"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/price_type"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/sales_representative"
)

// Handlers groups all partner-related handlers
type Handlers struct {
	Supplier       *supplier.Handler
	PartnerGroup   *partner_group.Handler
	PartnerSubGroup *partner_sub_group.Handler
	District       *district.Handler
	Thana          *thana.Handler
	TaxBracket     *tax_bracket.Handler
	Customer        *customer.Handler
	PriceType      *price_type.Handler
	SalesRepresentative *sales_representative.Handler
	
}
