package database

import (
	"fmt"
	"time"

	query_logger "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("mysql dsn is required")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

    sqlDB.SetMaxOpenConns(25)               
	sqlDB.SetMaxIdleConns(10)                
	sqlDB.SetConnMaxLifetime(5 * time.Minute) 

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}


	queryLogger := query_logger.NewQueryLogger(1000)
	db = queryLogger.GormCallback(db)
	return db, nil
}
