package main

import (
	"log"
	
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/app"
)

// @title ASSMI Super Shop ERP API
// @version 1.0
// @description Backend API for ASSMI Super Shop ERP.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
