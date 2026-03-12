package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"consulting-system/config"
	"consulting-system/internal/models"
	"consulting-system/internal/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	if err := initPostgreSQL(cfg); err != nil {
		return fmt.Errorf("failed to init PostgreSQL: %w", err)
	}

	if err := initRedis(cfg); err != nil {
		return fmt.Errorf("failed to init Redis: %w", err)
	}

	return nil
}

// initPostgreSQL 初始化PostgreSQL连接
func initPostgreSQL(cfg *config.Config) error {
	var logLevel logger.LogLevel
	switch cfg.Log.Level {
	case "debug":
		logLevel = logger.Info
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	default:
		logLevel = logger.Warn
	}

	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), dbConfig)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("PostgreSQL connected successfully")
	return nil
}

// initRedis 初始化Redis连接
func initRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Println("Redis connected successfully")
	return nil
}

// Migrate 执行数据库迁移
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Customer{},
		&models.FollowUpRecord{},
		&models.ServiceType{},
		&models.ServiceOrder{},
		&models.Communication{},
		&models.Invoice{},
		&models.PaymentPlan{},
		&models.Receipt{},
		&models.Refund{},
		&models.Notice{},
		&models.NoticeRead{},
		&models.Document{},
		&models.Contract{},
		&models.SystemConfig{},
		&models.OperationLog{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	log.Println("Database migration completed")
	return nil
}

// SeedData 初始化基础数据
func SeedData() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 检查是否已有角色数据
	var count int64
	if err := DB.Model(&models.Role{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Basic data already exists, skipping seed")
		return nil
	}

	// 创建默认角色
	roles := []models.Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "系统超级管理员，拥有所有权限",
			Permissions: []string{"*"},
			Status:      1,
		},
		{
			Name:        "管理员",
			Code:        "admin",
			Description: "系统管理员",
			Permissions: []string{
				"user:view", "user:create", "user:update",
				"customer:view", "customer:create", "customer:update", "customer:delete",
				"service:view", "service:create", "service:update", "service:delete",
				"invoice:view", "invoice:create", "invoice:update",
				"receipt:view", "receipt:create", "receipt:update",
				"statistics:view",
				"notice:view", "notice:create", "notice:update", "notice:delete",
				"document:view", "document:create", "document:update", "document:delete",
				"contract:view", "contract:create", "contract:update", "contract:delete",
				"setting:view", "setting:update",
			},
			Status: 1,
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户",
			Permissions: []string{
				"user:view",
				"customer:view",
				"service:view",
				"invoice:view",
				"receipt:view",
				"notice:view",
				"document:view",
				"contract:view",
			},
			Status: 1,
		},
	}

	for i := range roles {
		if err := DB.Create(&roles[i]).Error; err != nil {
			return fmt.Errorf("failed to create role %s: %w", roles[i].Name, err)
		}
	}

	// 创建默认管理员用户
	hashedPassword, _ := utils.HashPassword("admin123")

	adminUser := models.User{
		Username:   "admin",
		Password:   hashedPassword,
		Email:      "admin@example.com",
		Phone:      "13800138000",
		RealName:   "系统管理员",
		Department: "技术部",
		Status:     1,
	}

	if err := DB.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// 关联超级管理员角色
	if err := DB.Model(&adminUser).Association("Roles").Append(&roles[0]); err != nil {
		return fmt.Errorf("failed to assign role to admin: %w", err)
	}

	// 初始化系统配置
	configs := []models.SystemConfig{
		{Key: "site_name", Value: "咨询公司业务管理系统", Description: "站点名称"},
		{Key: "site_logo", Value: "", Description: "站点Logo"},
		{Key: "company_name", Value: "XX咨询公司", Description: "公司名称"},
		{Key: "contact_phone", Value: "400-000-0000", Description: "联系电话"},
		{Key: "contact_email", Value: "contact@example.com", Description: "联系邮箱"},
	}

	for i := range configs {
		if err := DB.Create(&configs[i]).Error; err != nil {
			return fmt.Errorf("failed to create config %s: %w", configs[i].Key, err)
		}
	}

	log.Println("Basic data seeded successfully")
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			return err
		}
	}

	return nil
}

// WithTransaction 执行事务
func WithTransaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}
