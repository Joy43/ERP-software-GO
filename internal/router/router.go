package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/config"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/middleware"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/audit"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/filter"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	rdb *redis.Client,
	cfg config.Config,
	h *Handlers,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(), filter.ExceptionHandler())
	//------locally and deploy url base adding Serve static files (media folder)
	r.Static("/media", "./media")
//---------base route group for versioning---------
	v1 := r.Group("/api/v1")
	// ----- Health check endpoint-----------
	v1.GET("/health", healthHandler)
	//----------------------------------------
	// db → MySQL/PostgreSQL database (GORM)
	// rdb → Redis client for caching and session management
	// h → Handlers struct containing service handlers for different modules (auth, audit, partner, etc.)
	// middleware.Auth(cfg) → Authentication middleware that checks for valid tokens and permissions based on the provided configuration
	//----------------------------------------
	auth.RegisterRoutes(v1, db, rdb, h.Auth, middleware.Auth(cfg))
	audit.RegisterRoutes(v1, db, rdb, h.Audit, middleware.Auth(cfg))
	partner.RegisterRoutes(v1, db, rdb, h.Partner, middleware.Auth(cfg))
	iteam_profile.RegisterRoutes(v1, db, rdb, h.ITeamProfile, middleware.Auth(cfg))
	uploads.RegisterRoutes(v1, db, rdb, h.Uploads, middleware.Auth(cfg))
	purchase.RegisterRoutes(v1, db, rdb, h.Purchase, middleware.Auth(cfg))

	//---------- Serve custom swagger UI with auth interceptor----------
	r.GET("/swagger/*any", func(c *gin.Context) {
		path := c.Param("any")
		switch path {
		case "/", "/index.html":
			c.File("./docs/swagger/index.html")
		case "/swagger-auth.js":
			c.File("./docs/swagger/swagger-auth.js")
		default:
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		}
	})

	return r
}

// healthHandler godoc
// @Summary Health check
// @Description Returns service health status
// @Tags System
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /health [get]
func healthHandler(ctx *gin.Context) {
	response.Success(ctx, "service is healthy", gin.H{"status": "ok"})
}
