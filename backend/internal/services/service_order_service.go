package services

import (
	"errors"
	"fmt"
	"time"

	"consulting-system/internal/database"
	"consulting-system/internal/models"
	"consulting-system/internal/utils"

	"gorm.io/gorm"
)

// ServiceOrderService 服务订单服务
type ServiceOrderService struct{}

// NewServiceOrderService 创建服务订单服务实例
func NewServiceOrderService() *ServiceOrderService {
	return &ServiceOrderService{}
}

// CreateServiceOrderRequest 创建服务订单请求
type CreateServiceOrderRequest struct {
	CustomerID    uint               `json:"customer_id" binding:"required"`
	ServiceTypeID uint               `json:"service_type_id" binding:"required"`
	Name          string             `json:"name" binding:"required,max=200"`
	Amount        float64            `json:"amount" binding:"omitempty,min=0"`
	StartDate     *time.Time         `json:"start_date" binding:"omitempty"`
	EndDate       *time.Time         `json:"end_date" binding:"omitempty"`
	Description   string             `json:"description" binding:"omitempty"`
	Participants  models.StringArray `json:"participants" binding:"omitempty"`
}

// UpdateServiceOrderRequest 更新服务订单请求
type UpdateServiceOrderRequest struct {
	Name         string             `json:"name" binding:"omitempty,max=200"`
	Amount       float64            `json:"amount" binding:"omitempty,min=0"`
	StartDate    *time.Time         `json:"start_date" binding:"omitempty"`
	EndDate      *time.Time         `json:"end_date" binding:"omitempty"`
	Status       int                `json:"status" binding:"omitempty,oneof=1 2 3 4"`
	Progress     int                `json:"progress" binding:"omitempty,min=0,max=100"`
	Description  string             `json:"description" binding:"omitempty"`
	Participants models.StringArray `json:"participants" binding:"omitempty"`
}

// UpdateProgressRequest 更新进度请求
type UpdateProgressRequest struct {
	Progress int `json:"progress" binding:"required,min=0,max=100"`
}

// CreateCommunicationRequest 创建沟通纪要请求
type CreateCommunicationRequest struct {
	Content string `json:"content" binding:"required"`
}

// ListServiceOrdersRequest 服务订单列表请求
type ListServiceOrdersRequest struct {
	CustomerID    uint   `form:"customer_id"`
	ServiceTypeID uint   `form:"service_type_id"`
	Name          string `form:"name"`
	Status        int    `form:"status"`
	Page          int    `form:"page,default=1"`
	PageSize      int    `form:"page_size,default=10"`
}

// ListServiceOrdersResponse 服务订单列表响应
type ListServiceOrdersResponse struct {
	List  []models.ServiceOrder `json:"list"`
	Total int64                 `json:"total"`
}

// generateServiceCode 生成服务单编号
func generateServiceCode() string {
	seq := int(time.Now().UnixNano() % 1000000)
	return utils.GenerateCode("SVC", seq)
}

// CreateServiceOrder 创建服务订单
func (s *ServiceOrderService) CreateServiceOrder(req *CreateServiceOrderRequest, createdBy uint) (*models.ServiceOrder, error) {
	// 检查客户是否存在
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		return nil, errors.New("客户不存在")
	}

	// 检查服务类型是否存在
	var serviceType models.ServiceType
	if err := database.DB.First(&serviceType, req.ServiceTypeID).Error; err != nil {
		return nil, errors.New("服务类型不存在")
	}

	order := models.ServiceOrder{
		CustomerID:    req.CustomerID,
		ServiceTypeID: req.ServiceTypeID,
		Code:          generateServiceCode(),
		Name:          req.Name,
		Amount:        req.Amount,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Status:        1, // 待启动
		Progress:      0,
		Description:   req.Description,
		Participants:  req.Participants,
		CreatedBy:     createdBy,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	database.DB.Preload("Customer").Preload("ServiceType").First(&order, order.ID)

	return &order, nil
}

// UpdateServiceOrder 更新服务订单
func (s *ServiceOrderService) UpdateServiceOrder(id uint, req *UpdateServiceOrderRequest) (*models.ServiceOrder, error) {
	var order models.ServiceOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务订单不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Amount >= 0 {
		updates["amount"] = req.Amount
	}
	if req.StartDate != nil {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = req.EndDate
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.Progress >= 0 {
		updates["progress"] = req.Progress
		// 如果进度为100%，自动更新状态为已完成
		if req.Progress == 100 {
			updates["status"] = 3
		}
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Participants != nil {
		updates["participants"] = req.Participants
	}

	if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").Preload("ServiceType").First(&order, id)

	return &order, nil
}

// UpdateProgress 更新服务进度
func (s *ServiceOrderService) UpdateProgress(id uint, req *UpdateProgressRequest) (*models.ServiceOrder, error) {
	var order models.ServiceOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务订单不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"progress": req.Progress,
	}

	// 如果进度为100%，自动更新状态为已完成
	if req.Progress == 100 {
		updates["status"] = 3
	} else if req.Progress > 0 && order.Status == 1 {
		// 如果进度大于0且状态为待启动，更新为进行中
		updates["status"] = 2
	}

	if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").Preload("ServiceType").First(&order, id)

	return &order, nil
}

// DeleteServiceOrder 删除服务订单
func (s *ServiceOrderService) DeleteServiceOrder(id uint) error {
	var order models.ServiceOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("服务订单不存在")
		}
		return err
	}

	return database.DB.Delete(&order).Error
}

// GetServiceOrder 获取服务订单详情
func (s *ServiceOrderService) GetServiceOrder(id uint) (*models.ServiceOrder, error) {
	var order models.ServiceOrder
	if err := database.DB.Preload("Customer").Preload("ServiceType").
		Preload("Communications", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(20)
		}).
		Preload("Communications.User").
		First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务订单不存在")
		}
		return nil, err
	}
	return &order, nil
}

// ListServiceOrders 获取服务订单列表
func (s *ServiceOrderService) ListServiceOrders(req *ListServiceOrdersRequest) (*ListServiceOrdersResponse, error) {
	var orders []models.ServiceOrder
	var total int64

	query := database.DB.Model(&models.ServiceOrder{})

	if req.CustomerID != 0 {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.ServiceTypeID != 0 {
		query = query.Where("service_type_id = ?", req.ServiceTypeID)
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Customer").Preload("ServiceType").
		Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	return &ListServiceOrdersResponse{
		List:  orders,
		Total: total,
	}, nil
}

// CreateCommunication 创建沟通纪要
func (s *ServiceOrderService) CreateCommunication(serviceID uint, req *CreateCommunicationRequest, createdBy uint) (*models.Communication, error) {
	// 检查服务订单是否存在
	var order models.ServiceOrder
	if err := database.DB.First(&order, serviceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务订单不存在")
		}
		return nil, err
	}

	communication := models.Communication{
		ServiceID: serviceID,
		Content:   req.Content,
		CreatedBy: createdBy,
	}

	if err := database.DB.Create(&communication).Error; err != nil {
		return nil, err
	}

	// 加载用户信息
	database.DB.Preload("User").First(&communication, communication.ID)

	return &communication, nil
}

// GetCommunications 获取沟通纪要列表
func (s *ServiceOrderService) GetCommunications(serviceID uint, page, pageSize int) ([]models.Communication, int64, error) {
	var communications []models.Communication
	var total int64

	query := database.DB.Model(&models.Communication{}).Where("service_id = ?", serviceID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&communications).Error; err != nil {
		return nil, 0, err
	}

	return communications, total, nil
}

// DeleteCommunication 删除沟通纪要
func (s *ServiceOrderService) DeleteCommunication(id uint) error {
	var communication models.Communication
	if err := database.DB.First(&communication, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("沟通纪要不存在")
		}
		return err
	}

	return database.DB.Delete(&communication).Error
}

// GetServiceStats 获取服务统计
func (s *ServiceOrderService) GetServiceStats() (map[string]interface{}, error) {
	var total int64
	var byStatus []struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}
	var totalAmount float64

	database.DB.Model(&models.ServiceOrder{}).Count(&total)
	database.DB.Model(&models.ServiceOrder{}).Select("status, COUNT(*) as count").Group("status").Scan(&byStatus)
	database.DB.Model(&models.ServiceOrder{}).Select("COALESCE(SUM(amount), 0) as total_amount").Pluck("total_amount", &totalAmount)

	return map[string]interface{}{
		"total":        total,
		"by_status":    byStatus,
		"total_amount": totalAmount,
	}, nil
}
