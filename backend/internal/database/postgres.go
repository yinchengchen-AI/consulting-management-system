package database

import (
	"consulting-system/backend/config"
	"consulting-system/backend/internal/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB 创建 PostgreSQL 数据库连接
func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// 配置 GORM 日志
	logLevel := logger.Silent
	if cfg.Host == "localhost" || cfg.Host == "127.0.0.1" {
		logLevel = logger.Info
	}

	// 打开数据库连接
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层 SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取 SQL DB 失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	if lifetime, err := time.ParseDuration(cfg.ConnMaxLifetime); err == nil {
		sqlDB.SetConnMaxLifetime(lifetime)
	}

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("自动迁移失败: %w", err)
	}

	return db, nil
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Project{},
		&models.Contract{},
		&models.Income{},
		&models.Expense{},
	)
}
