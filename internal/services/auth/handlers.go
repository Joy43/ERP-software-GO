package auth

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/department"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/designation"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/payment_mode"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/permission"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/responsibility_transfer"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/role"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/wallet"

)

// Handlers groups all auth-related handlers
type Handlers struct {
	User                   *user.Handler
	Role                   *role.Handler
	Permission             *permission.Handler
	Designation            *designation.Handler
	Department             *department.Handler
	Office                 *office.Handler
	PaymentMode            *payment_mode.Handler
	ResponsibilityTransfer *responsibility_transfer.Handler
	Wallet                 *wallet.Handler
	

}
