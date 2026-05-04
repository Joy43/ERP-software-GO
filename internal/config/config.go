package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName          string
	Port             string
	Env              string
	BaseURL          string
	MysqlDSN         string
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessTTL     time.Duration
	JWTRefreshTTL    time.Duration
	BCryptCost       int
	RedisHost        string
	RedisPort        string
	RedisPassword    string
}

func Load() (Config, error) {
	accessTTLMinutes, err := getEnvAsInt("JWT_ACCESS_TTL_MINUTES", 15)
	if err != nil {
		return Config{}, err
	}

	refreshTTLDays, err := getEnvAsInt("JWT_REFRESH_TTL_DAYS", 7)
	if err != nil {
		return Config{}, err
	}

	bcryptCost, err := getEnvAsInt("BCRYPT_COST", 12)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppName:          getEnv("APP_NAME", "assmi-super-shop-erp-backend"),
		Port:             getEnv("APP_PORT", "8080"),
		Env:              getEnv("APP_ENV", "development"),
		BaseURL:          getEnv("BASE_URL", "https://erp.vidatech.com.bd"),
		MysqlDSN:         getEnv("MYSQL_DSN", "root@tcp(localhost:3306)/assmi_super_shop?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
		JWTAccessTTL:     time.Duration(accessTTLMinutes) * time.Minute,
		JWTRefreshTTL:    time.Duration(refreshTTLDays) * 24 * time.Hour,
		BCryptCost:       bcryptCost,
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
	}

	if cfg.Port == "" {
		return Config{}, fmt.Errorf("APP_PORT is required")
	}

	if cfg.MysqlDSN == "" {
		return Config{}, fmt.Errorf("MYSQL_DSN is required")
	}

	if cfg.JWTAccessSecret == "" || cfg.JWTRefreshSecret == "" {
		return Config{}, fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET are required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	converted, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return converted, nil
}
