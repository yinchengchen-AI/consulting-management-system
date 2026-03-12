package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IncomeType 收入类型
type IncomeType string

const (
	IncomeTypeContract IncomeType = "contract" // 合同收入
	IncomeTypeProject  IncomeType = "project"  // 项目收入
	IncomeTypeOther    IncomeType = "other"    // 其他收入
)

// IncomeStatus 收入状态
type IncomeStatus string

const (
	IncomeStatusPending  IncomeStatus = "pending"  // 待确认
	IncomeStatusReceived IncomeStatus = "received" // 已收款
	IncomeStatusInvoiced IncomeStatus = "invoiced" // 已开票
)

// Income 收入记录模型
type Income struct {
	ID           string       `gorm:"type:uuid;primary_key" json:"id"`
	Type         IncomeType   `gorm:"type:varchar(20)" json:"type"`
	Status       IncomeStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Amount       float64      `gorm:"type:decimal(15,2);not null" json:"amount"`
	CustomerID   string       `gorm:"type:uuid" json:"customer_id"`
	Customer     *Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ContractID   string       `gorm:"type:uuid" json:"contract_id"`
	Contract     *Contract    `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	ProjectID    string       `gorm:"type:uuid" json:"project_id"`
	Project      *Project     `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Description  string       `gorm:"type:text" json:"description"`
	InvoiceNo    string       `gorm:"type:varchar(50)" json:"invoice_no"`
	InvoiceDate  *time.Time   `json:"invoice_date"`
	ReceivedDate *time.Time   `json:"received_date"`
	RecordedByID string       `gorm:"type:uuid" json:"recorded_by_id"`
	RecordedBy   *User        `gorm:"foreignKey:RecordedByID" json:"recorded_by,omitempty"`
	Notes        string       `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Income) TableName() string {
	return "incomes"
}

// BeforeCreate 创建前钩子
func (i *Income) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}

// ExpenseType 支出类型
type ExpenseType string

const (
	ExpenseTypeSalary    ExpenseType = "salary"    // 工资
	ExpenseTypeBonus     ExpenseType = "bonus"     // 奖金
	ExpenseTypeOffice    ExpenseType = "office"    // 办公费用
	ExpenseTypeTravel    ExpenseType = "travel"    // 差旅费
	ExpenseTypeMarketing ExpenseType = "marketing" // 市场费用
	ExpenseTypeProject   ExpenseType = "project"   // 项目成本
	ExpenseTypeTax       ExpenseType = "tax"       // 税费
	ExpenseTypeOther     ExpenseType = "other"     // 其他
)

// ExpenseStatus 支出状态
type ExpenseStatus string

const (
	ExpenseStatusPending   ExpenseStatus = "pending"   // 待支付
	ExpenseStatusPaid      ExpenseStatus = "paid"      // 已支付
	ExpenseStatusCancelled ExpenseStatus = "cancelled" // 已取消
)

// Expense 支出记录模型
type Expense struct {
	ID          string        `gorm:"type:uuid;primary_key" json:"id"`
	Type        ExpenseType   `gorm:"type:varchar(20)" json:"type"`
	Status      ExpenseStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Amount      float64       `gorm:"type:decimal(15,2);not null" json:"amount"`
	ProjectID   string        `gorm:"type:uuid" json:"project_id"`
	Project     *Project      `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Description string        `gorm:"type:text" json:"description"`
	Payee       string        `gorm:"type:varchar(100)" json:"payee"`
	PaidDate    *time.Time    `json:"paid_date"`
	RecordedByID string       `gorm:"type:uuid" json:"recorded_by_id"`
	RecordedBy  *User         `gorm:"foreignKey:RecordedByID" json:"recorded_by,omitempty"`
	Notes       string        `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Expense) TableName() string {
	return "expenses"
}

// BeforeCreate 创建前钩子
func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

// FinanceSummary 财务汇总
type FinanceSummary struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpense  float64 `json:"total_expense"`
	NetProfit     float64 `json:"net_profit"`
	ProfitMargin  float64 `json:"profit_margin"`
	PendingIncome float64 `json:"pending_income"`
	PendingExpense float64 `json:"pending_expense"`
}

// MonthlyFinance 月度财务数据
type MonthlyFinance struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Profit  float64 `json:"profit"`
}

// CreateIncomeRequest 创建收入请求
type CreateIncomeRequest struct {
	Type         IncomeType `json:"type" binding:"omitempty,oneof=contract project other"`
	Amount       float64    `json:"amount" binding:"required,gt=0"`
	CustomerID   string     `json:"customer_id"`
	ContractID   string     `json:"contract_id"`
	ProjectID    string     `json:"project_id"`
	Description  string     `json:"description"`
	InvoiceNo    string     `json:"invoice_no"`
	InvoiceDate  *time.Time `json:"invoice_date"`
	ReceivedDate *time.Time `json:"received_date"`
	Notes        string     `json:"notes"`
}

// CreateExpenseRequest 创建支出请求
type CreateExpenseRequest struct {
	Type        ExpenseType `json:"type" binding:"required,oneof=salary bonus office travel marketing project tax other"`
	Amount      float64     `json:"amount" binding:"required,gt=0"`
	ProjectID   string      `json:"project_id"`
	Description string      `json:"description"`
	Payee       string      `json:"payee"`
	PaidDate    *time.Time  `json:"paid_date"`
	Notes       string      `json:"notes"`
}
