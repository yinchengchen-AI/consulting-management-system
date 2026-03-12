package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectStatus 项目状态
type ProjectStatus string

const (
	ProjectStatusPending    ProjectStatus = "pending"    // 待启动
	ProjectStatusActive     ProjectStatus = "active"     // 进行中
	ProjectStatusPaused     ProjectStatus = "paused"     // 已暂停
	ProjectStatusCompleted  ProjectStatus = "completed"  // 已完成
	ProjectStatusCancelled  ProjectStatus = "cancelled"  // 已取消
)

// ProjectType 项目类型
type ProjectType string

const (
	ProjectTypeStrategy    ProjectType = "strategy"    // 战略咨询
	ProjectTypeManagement  ProjectType = "management"  // 管理咨询
	ProjectTypeTechnology  ProjectType = "technology"  // 技术咨询
	ProjectTypeFinance     ProjectType = "finance"     // 财务咨询
	ProjectTypeHR          ProjectType = "hr"          // 人力资源咨询
	ProjectTypeMarketing   ProjectType = "marketing"   // 市场营销咨询
	ProjectTypeOther       ProjectType = "other"       // 其他
)

// Project 项目模型
type Project struct {
	ID           string        `gorm:"type:uuid;primary_key" json:"id"`
	Name         string        `gorm:"type:varchar(100);not null" json:"name"`
	Code         string        `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	Type         ProjectType   `gorm:"type:varchar(20)" json:"type"`
	Status       ProjectStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Description  string        `gorm:"type:text" json:"description"`
	CustomerID   string        `gorm:"type:uuid;not null" json:"customer_id"`
	Customer     *Customer     `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ContractID   string        `gorm:"type:uuid" json:"contract_id"`
	Contract     *Contract     `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	ManagerID    string        `gorm:"type:uuid" json:"manager_id"`
	Manager      *User         `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	StartDate    *time.Time    `json:"start_date"`
	EndDate      *time.Time    `json:"end_date"`
	Budget       float64       `gorm:"type:decimal(15,2);default:0" json:"budget"`
	ActualCost   float64       `gorm:"type:decimal(15,2);default:0" json:"actual_cost"`
	Progress     int           `gorm:"default:0" json:"progress"`
	Priority     int           `gorm:"default:3" json:"priority"` // 1-最高, 5-最低
	Deliverables string        `gorm:"type:text" json:"deliverables"`
	Notes        string        `gorm:"type:text" json:"notes"`
	Members      []ProjectMember `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
	Tasks        []ProjectTask   `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate 创建前钩子
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.Code == "" {
		p.Code = generateProjectCode()
	}
	return nil
}

// generateProjectCode 生成项目编号
func generateProjectCode() string {
	return "PRJ-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
}

// ProjectMember 项目成员模型
type ProjectMember struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID string    `gorm:"type:uuid;not null" json:"project_id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      string    `gorm:"type:varchar(50)" json:"role"`
	JoinDate  time.Time `json:"join_date"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (ProjectMember) TableName() string {
	return "project_members"
}

// BeforeCreate 创建前钩子
func (pm *ProjectMember) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == "" {
		pm.ID = uuid.New().String()
	}
	return nil
}

// ProjectTask 项目任务模型
type ProjectTask struct {
	ID          string     `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID   string     `gorm:"type:uuid;not null" json:"project_id"`
	Title       string     `gorm:"type:varchar(100);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	AssigneeID  string     `gorm:"type:uuid" json:"assignee_id"`
	Assignee    *User      `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Status      string     `gorm:"type:varchar(20);default:'todo'" json:"status"`
	Priority    int        `gorm:"default:3" json:"priority"`
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (ProjectTask) TableName() string {
	return "project_tasks"
}

// BeforeCreate 创建前钩子
func (pt *ProjectTask) BeforeCreate(tx *gorm.DB) error {
	if pt.ID == "" {
		pt.ID = uuid.New().String()
	}
	return nil
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name        string      `json:"name" binding:"required,max=100"`
	Type        ProjectType `json:"type" binding:"omitempty,oneof=strategy management technology finance hr marketing other"`
	Description string      `json:"description"`
	CustomerID  string      `json:"customer_id" binding:"required,uuid"`
	ContractID  string      `json:"contract_id"`
	ManagerID   string      `json:"manager_id"`
	StartDate   *time.Time  `json:"start_date"`
	EndDate     *time.Time  `json:"end_date"`
	Budget      float64     `json:"budget"`
	Priority    int         `json:"priority" binding:"omitempty,min=1,max=5"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Name        string        `json:"name" binding:"omitempty,max=100"`
	Type        ProjectType   `json:"type" binding:"omitempty,oneof=strategy management technology finance hr marketing other"`
	Status      ProjectStatus `json:"status" binding:"omitempty,oneof=pending active paused completed cancelled"`
	Description string        `json:"description"`
	ManagerID   string        `json:"manager_id"`
	StartDate   *time.Time    `json:"start_date"`
	EndDate     *time.Time    `json:"end_date"`
	Budget      float64       `json:"budget"`
	ActualCost  float64       `json:"actual_cost"`
	Progress    int           `json:"progress" binding:"omitempty,min=0,max=100"`
	Priority    int           `json:"priority" binding:"omitempty,min=1,max=5"`
	Deliverables string       `json:"deliverables"`
	Notes       string        `json:"notes"`
}
