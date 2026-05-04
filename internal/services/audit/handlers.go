package audit

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/audit/transaction_history"
)

// Handlers groups all audit-related handlers
type Handlers struct {
	TransactionHistory *transaction_history.Handler
}
