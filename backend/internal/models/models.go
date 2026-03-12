package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSON json类型
type JSON map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &j)
}

// StringArray 字符串数组类型
type StringArray []string

// Value 实现driver.Valuer接口
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 实现sql.Scanner接口
func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &a)
}

// ==================== 1. 用户权限管理模块 ====================

// User 用户模型
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Email       string         `gorm:"uniqueIndex;size:100" json:"email"`
	Phone       string         `gorm:"size:20" json:"phone"`
	RealName    string         `gorm:"size:50" json:"real_name"`
	Department  string         `gorm:"size:50" json:"department"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Status      int            `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Roles       []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	Permissions StringArray    `gorm:"type:jsonb" json:"permissions"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Users       []User         `gorm:"many2many:user_roles;" json:"-"`
}

// UserRole 用户角色关联
type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== 2. 客户公司管理模块 ====================

// Customer 客户公司模型
type Customer struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Industry      string         `gorm:"size:50" json:"industry"`
	Scale         string         `gorm:"size:20;comment:企业规模" json:"scale"`
	ContactName   string         `gorm:"size:50" json:"contact_name"`
	ContactPhone  string         `gorm:"size:20" json:"contact_phone"`
	ContactEmail  string         `gorm:"size:100" json:"contact_email"`
	Address       string         `gorm:"size:255" json:"address"`
	Website       string         `gorm:"size:100" json:"website"`
	Status        int            `gorm:"default:1;comment:1-潜在客户,2-合作中,3-暂停,4-终止" json:"status"`
	Tags          StringArray    `gorm:"type:jsonb" json:"tags"`
	Remark        string         `gorm:"type:text" json:"remark"`
	CreatedBy     uint           `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	FollowUpRecords []FollowUpRecord `json:"follow_up_records,omitempty"`
}

// FollowUpRecord 客户跟进记录
type FollowUpRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CustomerID uint      `gorm:"not null;index" json:"customer_id"`
	FollowUpBy uint      `gorm:"not null" json:"follow_up_by"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	FollowUpAt time.Time `json:"follow_up_at"`
	CreatedAt  time.Time `json:"created_at"`
	Customer   Customer  `json:"-"`
	User       User      `gorm:"foreignKey:FollowUpBy" json:"user,omitempty"`
}

// ==================== 3. 服务类型管理模块 ====================

// ServiceType 服务类型模型
type ServiceType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	Level       int            `gorm:"default:1" json:"level"`
	Path        string         `gorm:"size:255" json:"path"`
	PriceMin    float64        `gorm:"type:decimal(15,2)" json:"price_min"`
	PriceMax    float64        `gorm:"type:decimal(15,2)" json:"price_max"`
	TaxRate     float64        `gorm:"type:decimal(5,2);default:6" json:"tax_rate"`
	Template    JSON           `gorm:"type:jsonb" json:"template"`
	Description string         `gorm:"type:text" json:"description"`
	Status      int            `gorm:"default:1" json:"status"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Parent      *ServiceType   `json:"parent,omitempty"`
	Children    []ServiceType  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// ==================== 4. 服务详细内容管理模块 ====================

// ServiceOrder 服务单模型
type ServiceOrder struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CustomerID   uint           `gorm:"not null;index" json:"customer_id"`
	ServiceTypeID uint          `gorm:"not null" json:"service_type_id"`
	Code         string         `gorm:"uniqueIndex;size:50" json:"code"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	Amount       float64        `gorm:"type:decimal(15,2)" json:"amount"`
	StartDate    *time.Time     `json:"start_date"`
	EndDate      *time.Time     `json:"end_date"`
	Status       int            `gorm:"default:1;comment:1-待启动,2-进行中,3-已完成,4-已暂停" json:"status"`
	Progress     int            `gorm:"default:0;comment:进度百分比" json:"progress"`
	Description  string         `gorm:"type:text" json:"description"`
	Participants StringArray    `gorm:"type:jsonb" json:"participants"`
	CreatedBy    uint           `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Customer     Customer       `json:"customer,omitempty"`
	ServiceType  ServiceType    `json:"service_type,omitempty"`
	Communications []Communication `json:"communications,omitempty"`
}

// Communication 沟通纪要
type Communication struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	ServiceID uint         `gorm:"not null;index" json:"service_id"`
	Content   string       `gorm:"type:text;not null" json:"content"`
	CreatedBy uint         `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	Service   ServiceOrder `json:"-"`
	User      User         `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

// ==================== 5. 开票信息管理模块 ====================

// Invoice 开票记录模型
type Invoice struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CustomerID   uint           `gorm:"not null;index" json:"customer_id"`
	ServiceID    *uint          `json:"service_id"`
	InvoiceType  int            `gorm:"default:1;comment:1-普票,2-专票" json:"invoice_type"`
	Amount       float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	TaxRate      float64        `gorm:"type:decimal(5,2)" json:"tax_rate"`
	TaxAmount    float64        `gorm:"type:decimal(15,2)" json:"tax_amount"`
	TotalAmount  float64        `gorm:"type:decimal(15,2)" json:"total_amount"`
	InvoiceNo    string         `gorm:"size:50" json:"invoice_no"`
	InvoiceCode  string         `gorm:"size:50" json:"invoice_code"`
	Status       int            `gorm:"default:1;comment:1-待开票,2-已开票,3-已作废" json:"status"`
	InvoiceInfo  JSON           `gorm:"type:jsonb" json:"invoice_info"`
	InvoiceDate  *time.Time     `json:"invoice_date"`
	CreatedBy    uint           `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Customer     Customer       `json:"customer,omitempty"`
}

// ==================== 6. 收款信息管理模块 ====================

// PaymentPlan 收款计划模型
type PaymentPlan struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CustomerID  uint           `gorm:"not null;index" json:"customer_id"`
	ServiceID   *uint          `json:"service_id"`
	InvoiceID   *uint          `json:"invoice_id"`
	Amount      float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	PlannedDate *time.Time     `json:"planned_date"`
	Status      int            `gorm:"default:1;comment:1-待收款,2-部分收款,3-已收完" json:"status"`
	Remark      string         `gorm:"type:text" json:"remark"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Customer    Customer       `json:"customer,omitempty"`
	Receipts    []Receipt      `json:"receipts,omitempty"`
}

// Receipt 收款记录模型
type Receipt struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PlanID        uint           `gorm:"not null;index" json:"plan_id"`
	Amount        float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	ReceivedDate  time.Time      `json:"received_date"`
	PaymentMethod int            `gorm:"default:1;comment:1-银行转账,2-现金,3-支票,4-其他" json:"payment_method"`
	Account       string         `gorm:"size:100" json:"account"`
	Remark        string         `gorm:"type:text" json:"remark"`
	CreatedBy     uint           `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	Plan          PaymentPlan    `json:"-"`
	Refunds       []Refund       `json:"refunds,omitempty"`
}

// Refund 退款记录模型
type Refund struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	ReceiptID uint     `gorm:"not null" json:"receipt_id"`
	Amount   float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Reason   string    `gorm:"type:text" json:"reason"`
	Status   int       `gorm:"default:1;comment:1-待审核,2-已批准,3-已拒绝,4-已完成" json:"status"`
	CreatedBy uint     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	Receipt  Receipt   `json:"-"`
}

// ==================== 8. 通知公告模块 ====================

// Notice 通知公告模型
type Notice struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Type        int            `gorm:"default:1;comment:1-普通,2-紧急,3-系统" json:"type"`
	TargetRoles StringArray    `gorm:"type:jsonb" json:"target_roles"`
	IsTop       bool           `gorm:"default:false" json:"is_top"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	ReadRecords []NoticeRead   `json:"read_records,omitempty"`
}

// NoticeRead 通知阅读记录
type NoticeRead struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	NoticeID uint      `gorm:"not null;index" json:"notice_id"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	ReadAt   time.Time `json:"read_at"`
	Notice   Notice    `json:"-"`
	User     User      `json:"-"`
}

// ==================== 9. 文档管理模块 ====================

// Document 文档模型
type Document struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"size:255;not null" json:"name"`
	Type             string         `gorm:"size:50" json:"type"`
	FileURL          string         `gorm:"size:500" json:"file_url"`
	Size             int64          `json:"size"`
	MimeType         string         `gorm:"size:100" json:"mime_type"`
	RelatedType      string         `gorm:"size:50;comment:关联类型:customer,service,contract等" json:"related_type"`
	RelatedID        uint           `json:"related_id"`
	AccessPermissions JSON          `gorm:"type:jsonb" json:"access_permissions"`
	Description      string         `gorm:"type:text" json:"description"`
	CreatedBy        uint           `json:"created_by"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	User             User           `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

// ==================== 10. 合同管理模块 ====================

// Contract 合同模型
type Contract struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CustomerID    uint           `gorm:"not null;index" json:"customer_id"`
	ServiceID     *uint          `json:"service_id"`
	Code          string         `gorm:"uniqueIndex;size:50" json:"code"`
	Name          string         `gorm:"size:200;not null" json:"name"`
	Amount        float64        `gorm:"type:decimal(15,2)" json:"amount"`
	SignDate      *time.Time     `json:"sign_date"`
	ExpireDate    *time.Time     `json:"expire_date"`
	PaymentTerms  string         `gorm:"type:text" json:"payment_terms"`
	Status        int            `gorm:"default:1;comment:1-草稿,2-已签署,3-履行中,4-已完成,5-已终止" json:"status"`
	FileURL       string         `gorm:"size:500" json:"file_url"`
	Remark        string         `gorm:"type:text" json:"remark"`
	CreatedBy     uint           `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Customer      Customer       `json:"customer,omitempty"`
}

// ==================== 11. 系统设置模块 ====================

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value       string    `gorm:"type:text" json:"value"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:100;not null" json:"action"`
	TargetType string    `gorm:"size:50" json:"target_type"`
	TargetID   uint      `json:"target_id"`
	Details    JSON      `gorm:"type:jsonb" json:"details"`
	IP         string    `gorm:"size:50" json:"ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
