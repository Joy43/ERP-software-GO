package router

import (
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/audit"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads"
)

// --- Handlers aggregates all domain-specific handlers for the application------
type Handlers struct {
	Auth        *auth.Handlers
	Partner     *partner.Handlers
	Audit       *audit.Handlers
	ITeamProfile *iteam_profile.Handlers
	Uploads     *uploads.Handler
	Purchase    *purchase.Handlers
	
}

