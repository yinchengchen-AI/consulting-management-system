package services

import (
	"errors"
	"time"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"gorm.io/gorm"
)

// CustomerService 客户服务
type CustomerService struct{}

// NewCustomerService 创建客户服务实例
func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

// CreateCustomerRequest 创建客户请求
type CreateCustomerRequest struct {
	Name         string              `json:"name" binding:"required,max=100"`
	Industry     string              `json:"industry" binding:"omitempty,max=50"`
	Scale        string              `json:"scale" binding:"omitempty,max=20"`
	ContactName  string              `json:"contact_name" binding:"omitempty,max=50"`
	ContactPhone string              `json:"contact_phone" binding:"omitempty,max=20"`
	ContactEmail string              `json:"contact_email" binding:"omitempty,email,max=100"`
	Address      string              `json:"address" binding:"omitempty,max=255"`
	Website      string              `json:"website" binding:"omitempty,max=100"`
	Tags         models.StringArray  `json:"tags" binding:"omitempty"`
	Remark       string              `json:"remark" binding:"omitempty"`
	Status       int                 `json:"status" binding:"omitempty,oneof=1 2 3 4"`
}

// UpdateCustomerRequest 更新客户请求
type UpdateCustomerRequest struct {
	Name         string              `json:"name" binding:"omitempty,max=100"`
	Industry     string              `json:"industry" binding:"omitempty,max=50"`
	Scale        string              `json:"scale" binding:"omitempty,max=20"`
	ContactName  string              `json:"contact_name" binding:"omitempty,max=50"`
	ContactPhone string              `json:"contact_phone" binding:"omitempty,max=20"`
	ContactEmail string              `json:"contact_email" binding:"omitempty,email,max=100"`
	Address      string              `json:"address" binding:"omitempty,max=255"`
	Website      string              `json:"website" binding:"omitempty,max=100"`
	Tags         models.StringArray  `json:"tags" binding:"omitempty"`
	Remark       string              `json:"remark" binding:"omitempty"`
	Status       int                 `json:"status" binding:"omitempty,oneof=1 2 3 4"`
}

// ListCustomersRequest 客户列表请求
type ListCustomersRequest struct {
	Name     string `form:"name"`
	Industry string `form:"industry"`
	Scale    string `form:"scale"`
	Status   int    `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

// ListCustomersResponse 客户列表响应
type ListCustomersResponse struct {
	List  []models.Customer `json:"list"`
	Total int64             `json:"total"`
}

// CreateFollowUpRequest 创建跟进记录请求
type CreateFollowUpRequest struct {
	Content    string    `json:"content" binding:"required"`
	FollowUpAt time.Time `json:"follow_up_at" binding:"omitempty"`
}

// CreateCustomer 创建客户
func (s *CustomerService) CreateCustomer(req *CreateCustomerRequest, createdBy uint) (*models.Customer, error) {
	customer := models.Customer{
		Name:         req.Name,
		Industry:     req.Industry,
		Scale:        req.Scale,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		Address:      req.Address,
		Website:      req.Website,
		Tags:         req.Tags,
		Remark:       req.Remark,
		Status:       req.Status,
		CreatedBy:    createdBy,
	}

	if customer.Status == 0 {
		customer.Status = 1
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

// UpdateCustomer 更新客户
func (s *CustomerService) UpdateCustomer(id uint, req *UpdateCustomerRequest) (*models.Customer, error) {
	var customer models.Customer
	if err := database.DB.First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("客户不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Industry != "" {
		updates["industry"] = req.Industry
	}
	if req.Scale != "" {
		updates["scale"] = req.Scale
	}
	if req.ContactName != "" {
		updates["contact_name"] = req.ContactName
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.ContactEmail != "" {
		updates["contact_email"] = req.ContactEmail
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := database.DB.Model(&customer).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

// DeleteCustomer 删除客户
func (s *CustomerService) DeleteCustomer(id uint) error {
	var customer models.Customer
	if err := database.DB.First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("客户不存在")
		}
		return err
	}

	return database.DB.Delete(&customer).Error
}

// GetCustomer 获取客户详情
func (s *CustomerService) GetCustomer(id uint) (*models.Customer, error) {
	var customer models.Customer
	if err := database.DB.Preload("FollowUpRecords", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(10)
	}).Preload("FollowUpRecords.User").First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("客户不存在")
		}
		return nil, err
	}
	return &customer, nil
}

// ListCustomers 获取客户列表
func (s *CustomerService) ListCustomers(req *ListCustomersRequest) (*ListCustomersResponse, error) {
	var customers []models.Customer
	var total int64

	query := database.DB.Model(&models.Customer{})

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Industry != "" {
		query = query.Where("industry = ?", req.Industry)
	}
	if req.Scale != "" {
		query = query.Where("scale = ?", req.Scale)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, err
	}

	return &ListCustomersResponse{
		List:  customers,
		Total: total,
	}, nil
}

// CreateFollowUp 创建跟进记录
func (s *CustomerService) CreateFollowUp(customerID uint, req *CreateFollowUpRequest, createdBy uint) (*models.FollowUpRecord, error) {
	// 检查客户是否存在
	var customer models.Customer
	if err := database.DB.First(&customer, customerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("客户不存在")
		}
		return nil, err
	}

	followUp := models.FollowUpRecord{
		CustomerID: customerID,
		FollowUpBy: createdBy,
		Content:    req.Content,
		FollowUpAt: req.FollowUpAt,
	}

	if followUp.FollowUpAt.IsZero() {
		followUp.FollowUpAt = time.Now()
	}

	if err := database.DB.Create(&followUp).Error; err != nil {
		return nil, err
	}

	// 加载用户信息
	database.DB.Preload("User").First(&followUp, followUp.ID)

	return &followUp, nil
}

// GetFollowUpRecords 获取跟进记录列表
func (s *CustomerService) GetFollowUpRecords(customerID uint, page, pageSize int) ([]models.FollowUpRecord, int64, error) {
	var records []models.FollowUpRecord
	var total int64

	query := database.DB.Model(&models.FollowUpRecord{}).Where("customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// DeleteFollowUp 删除跟进记录
func (s *CustomerService) DeleteFollowUp(id uint) error {
	var record models.FollowUpRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("跟进记录不存在")
		}
		return err
	}

	return database.DB.Delete(&record).Error
}

// GetCustomerStats 获取客户统计
func (s *CustomerService) GetCustomerStats() (map[string]interface{}, error) {
	var total int64
	var byStatus []struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}
	var byIndustry []struct {
		Industry string `json:"industry"`
		Count    int64  `json:"count"`
	}

	database.DB.Model(&models.Customer{}).Count(&total)
	database.DB.Model(&models.Customer{}).Select("status, COUNT(*) as count").Group("status").Scan(&byStatus)
	database.DB.Model(&models.Customer{}).Select("industry, COUNT(*) as count").Where("industry != ?", "").Group("industry").Scan(&byIndustry)

	return map[string]interface{}{
		"total":       total,
		"by_status":   byStatus,
		"by_industry": byIndustry,
	}, nil
}
