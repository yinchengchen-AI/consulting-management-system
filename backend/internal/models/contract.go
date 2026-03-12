package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContractStatus 合同状态
type ContractStatus string

const (
	ContractStatusDraft     ContractStatus = "draft"     // 草稿
	ContractStatusPending   ContractStatus = "pending"   // 待审批
	ContractStatusActive    ContractStatus = "active"    // 生效中
	ContractStatusCompleted ContractStatus = "completed" // 已完成
	ContractStatusCancelled ContractStatus = "cancelled" // 已取消
	ContractStatusExpired   ContractStatus = "expired"   // 已过期
)

// ContractType 合同类型
type ContractType string

const (
	ContractTypeProject   ContractType = "project"   // 项目合同
	ContractTypeRetainer  ContractType = "retainer"  // 年度框架合同
	ContractTypeConsulting ContractType = "consulting" // 单次咨询合同
)

// PaymentTerms 付款条款
type PaymentTerms string

const (
	PaymentTermsPrepay    PaymentTerms = "prepay"    // 预付
	PaymentTermsMilestone PaymentTerms = "milestone" // 里程碑付款
	PaymentTermsMonthly   PaymentTerms = "monthly"   // 月付
	PaymentTermsPostpay   PaymentTerms = "postpay"   // 后付
)

// Contract 合同模型
type Contract struct {
	ID             string         `gorm:"type:uuid;primary_key" json:"id"`
	Code           string         `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Type           ContractType   `gorm:"type:varchar(20)" json:"type"`
	Status         ContractStatus `gorm:"type:varchar(20);default:'draft'" json:"status"`
	CustomerID     string         `gorm:"type:uuid;not null" json:"customer_id"`
	Customer       *Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Amount         float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	TaxRate        float64        `gorm:"type:decimal(5,2);default:6" json:"tax_rate"`
	TaxAmount      float64        `gorm:"type:decimal(15,2)" json:"tax_amount"`
	TotalAmount    float64        `gorm:"type:decimal(15,2)" json:"total_amount"`
	PaymentTerms   PaymentTerms   `gorm:"type:varchar(20)" json:"payment_terms"`
	SignedDate     *time.Time     `json:"signed_date"`
	StartDate      *time.Time     `json:"start_date"`
	EndDate        *time.Time     `json:"end_date"`
	Description    string         `gorm:"type:text" json:"description"`
	Terms          string         `gorm:"type:text" json:"terms"`
	SalesOwnerID   string         `gorm:"type:uuid" json:"sales_owner_id"`
	SalesOwner     *User          `gorm:"foreignKey:SalesOwnerID" json:"sales_owner,omitempty"`
	SignedBy       string         `gorm:"type:varchar(50)" json:"signed_by"`
	AttachmentURL  string         `gorm:"type:varchar(255)" json:"attachment_url"`
	PaidAmount     float64        `gorm:"type:decimal(15,2);default:0" json:"paid_amount"`
	RemainingAmount float64       `gorm:"type:decimal(15,2);default:0" json:"remaining_amount"`
	Projects       []Project      `gorm:"foreignKey:ContractID" json:"projects,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Contract) TableName() string {
	return "contracts"
}

// BeforeCreate 创建前钩子
func (c *Contract) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.Code == "" {
		c.Code = generateContractCode()
	}
	// 计算税额和总金额
	c.TaxAmount = c.Amount * c.TaxRate / 100
	c.TotalAmount = c.Amount + c.TaxAmount
	c.RemainingAmount = c.TotalAmount
	return nil
}

// BeforeUpdate 更新前钩子
func (c *Contract) BeforeUpdate(tx *gorm.DB) error {
	// 重新计算税额和总金额
	c.TaxAmount = c.Amount * c.TaxRate / 100
	c.TotalAmount = c.Amount + c.TaxAmount
	c.RemainingAmount = c.TotalAmount - c.PaidAmount
	return nil
}

// generateContractCode 生成合同编号
func generateContractCode() string {
	return "CNT-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
}

// CreateContractRequest 创建合同请求
type CreateContractRequest struct {
	Name         string       `json:"name" binding:"required,max=100"`
	Type         ContractType `json:"type" binding:"omitempty,oneof=project retainer consulting"`
	CustomerID   string       `json:"customer_id" binding:"required,uuid"`
	Amount       float64      `json:"amount" binding:"required,gt=0"`
	TaxRate      float64      `json:"tax_rate" binding:"omitempty,min=0,max=100"`
	PaymentTerms PaymentTerms `json:"payment_terms" binding:"omitempty,oneof=prepay milestone monthly postpay"`
	SignedDate   *time.Time   `json:"signed_date"`
	StartDate    *time.Time   `json:"start_date"`
	EndDate      *time.Time   `json:"end_date"`
	Description  string       `json:"description"`
	Terms        string       `json:"terms"`
	SalesOwnerID string       `json:"sales_owner_id"`
	SignedBy     string       `json:"signed_by"`
}

// UpdateContractRequest 更新合同请求
type UpdateContractRequest struct {
	Name         string         `json:"name" binding:"omitempty,max=100"`
	Type         ContractType   `json:"type" binding:"omitempty,oneof=project retainer consulting"`
	Status       ContractStatus `json:"status" binding:"omitempty,oneof=draft pending active completed cancelled expired"`
	Amount       float64        `json:"amount" binding:"omitempty,gt=0"`
	TaxRate      float64        `json:"tax_rate" binding:"omitempty,min=0,max=100"`
	PaymentTerms PaymentTerms   `json:"payment_terms" binding:"omitempty,oneof=prepay milestone monthly postpay"`
	SignedDate   *time.Time     `json:"signed_date"`
	StartDate    *time.Time     `json:"start_date"`
	EndDate      *time.Time     `json:"end_date"`
	Description  string         `json:"description"`
	Terms        string         `json:"terms"`
	SalesOwnerID string         `json:"sales_owner_id"`
	SignedBy     string         `json:"signed_by"`
}
