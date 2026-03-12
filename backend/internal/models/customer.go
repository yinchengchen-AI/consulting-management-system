package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomerType 客户类型
type CustomerType string

const (
	CustomerTypeEnterprise CustomerType = "enterprise"
	CustomerTypeGovernment CustomerType = "government"
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeOther      CustomerType = "other"
)

// CustomerLevel 客户等级
type CustomerLevel string

const (
	CustomerLevelA CustomerLevel = "A" // 战略客户
	CustomerLevelB CustomerLevel = "B" // 重要客户
	CustomerLevelC CustomerLevel = "C" // 普通客户
	CustomerLevelD CustomerLevel = "D" // 潜在客户
)

// CustomerStatus 客户状态
type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusInactive  CustomerStatus = "inactive"
	CustomerStatusPotential CustomerStatus = "potential"
)

// Customer 客户模型
type Customer struct {
	ID            string         `gorm:"type:uuid;primary_key" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null" json:"name"`
	Type          CustomerType   `gorm:"type:varchar(20);default:'enterprise'" json:"type"`
	Level         CustomerLevel  `gorm:"type:varchar(10);default:'C'" json:"level"`
	Status        CustomerStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	Industry      string         `gorm:"type:varchar(50)" json:"industry"`
	Scale         string         `gorm:"type:varchar(50)" json:"scale"`
	Website       string         `gorm:"type:varchar(255)" json:"website"`
	Address       string         `gorm:"type:varchar(255)" json:"address"`
	Description   string         `gorm:"type:text" json:"description"`
	ContactName   string         `gorm:"type:varchar(50)" json:"contact_name"`
	ContactPhone  string         `gorm:"type:varchar(20)" json:"contact_phone"`
	ContactEmail  string         `gorm:"type:varchar(100)" json:"contact_email"`
	ContactTitle  string         `gorm:"type:varchar(50)" json:"contact_title"`
	SalesOwnerID  string         `gorm:"type:uuid" json:"sales_owner_id"`
	SalesOwner    *User          `gorm:"foreignKey:SalesOwnerID" json:"sales_owner,omitempty"`
	ContractCount int            `gorm:"-" json:"contract_count"`
	ProjectCount  int            `gorm:"-" json:"project_count"`
	TotalAmount   float64        `gorm:"-" json:"total_amount"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Customer) TableName() string {
	return "customers"
}

// BeforeCreate 创建前钩子
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// CreateCustomerRequest 创建客户请求
type CreateCustomerRequest struct {
	Name         string       `json:"name" binding:"required,max=100"`
	Type         CustomerType `json:"type" binding:"omitempty,oneof=enterprise government individual other"`
	Level        CustomerLevel `json:"level" binding:"omitempty,oneof=A B C D"`
	Industry     string       `json:"industry"`
	Scale        string       `json:"scale"`
	Website      string       `json:"website"`
	Address      string       `json:"address"`
	Description  string       `json:"description"`
	ContactName  string       `json:"contact_name"`
	ContactPhone string       `json:"contact_phone"`
	ContactEmail string       `json:"contact_email,email"`
	ContactTitle string       `json:"contact_title"`
	SalesOwnerID string       `json:"sales_owner_id"`
}

// UpdateCustomerRequest 更新客户请求
type UpdateCustomerRequest struct {
	Name         string         `json:"name" binding:"omitempty,max=100"`
	Type         CustomerType   `json:"type" binding:"omitempty,oneof=enterprise government individual other"`
	Level        CustomerLevel  `json:"level" binding:"omitempty,oneof=A B C D"`
	Status       CustomerStatus `json:"status" binding:"omitempty,oneof=active inactive potential"`
	Industry     string         `json:"industry"`
	Scale        string         `json:"scale"`
	Website      string         `json:"website"`
	Address      string         `json:"address"`
	Description  string         `json:"description"`
	ContactName  string         `json:"contact_name"`
	ContactPhone string         `json:"contact_phone"`
	ContactEmail string         `json:"contact_email,email"`
	ContactTitle string         `json:"contact_title"`
	SalesOwnerID string         `json:"sales_owner_id"`
}
