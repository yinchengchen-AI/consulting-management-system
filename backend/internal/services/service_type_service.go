package services

import (
	"errors"
	"fmt"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"gorm.io/gorm"
)

// ServiceTypeService 服务类型服务
type ServiceTypeService struct{}

// NewServiceTypeService 创建服务类型服务实例
func NewServiceTypeService() *ServiceTypeService {
	return &ServiceTypeService{}
}

// CreateServiceTypeRequest 创建服务类型请求
type CreateServiceTypeRequest struct {
	Name        string      `json:"name" binding:"required,max=100"`
	Code        string      `json:"code" binding:"required,max=50"`
	ParentID    *uint       `json:"parent_id" binding:"omitempty"`
	PriceMin    float64     `json:"price_min" binding:"omitempty,min=0"`
	PriceMax    float64     `json:"price_max" binding:"omitempty,min=0"`
	TaxRate     float64     `json:"tax_rate" binding:"omitempty,min=0,max=100"`
	Template    models.JSON `json:"template" binding:"omitempty"`
	Description string      `json:"description" binding:"omitempty"`
	Status      int         `json:"status" binding:"omitempty,oneof=1 2"`
	SortOrder   int         `json:"sort_order" binding:"omitempty"`
}

// UpdateServiceTypeRequest 更新服务类型请求
type UpdateServiceTypeRequest struct {
	Name        string      `json:"name" binding:"omitempty,max=100"`
	PriceMin    float64     `json:"price_min" binding:"omitempty,min=0"`
	PriceMax    float64     `json:"price_max" binding:"omitempty,min=0"`
	TaxRate     float64     `json:"tax_rate" binding:"omitempty,min=0,max=100"`
	Template    models.JSON `json:"template" binding:"omitempty"`
	Description string      `json:"description" binding:"omitempty"`
	Status      int         `json:"status" binding:"omitempty,oneof=1 2"`
	SortOrder   int         `json:"sort_order" binding:"omitempty"`
}

// ListServiceTypesRequest 服务类型列表请求
type ListServiceTypesRequest struct {
	Name     string `form:"name"`
	Code     string `form:"code"`
	ParentID *uint  `form:"parent_id"`
	Status   int    `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

// ListServiceTypesResponse 服务类型列表响应
type ListServiceTypesResponse struct {
	List  []models.ServiceType `json:"list"`
	Total int64                `json:"total"`
}

// CreateServiceType 创建服务类型
func (s *ServiceTypeService) CreateServiceType(req *CreateServiceTypeRequest) (*models.ServiceType, error) {
	// 检查代码是否已存在
	var existing models.ServiceType
	if err := database.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("服务类型代码已存在")
	}

	// 如果有父级，检查父级是否存在
	level := 1
	path := ""
	if req.ParentID != nil {
		var parent models.ServiceType
		if err := database.DB.First(&parent, *req.ParentID).Error; err != nil {
			return nil, errors.New("父级服务类型不存在")
		}
		level = parent.Level + 1
		path = parent.Path + "/" + fmt.Sprintf("%d", parent.ID)
	}

	serviceType := models.ServiceType{
		Name:        req.Name,
		Code:        req.Code,
		ParentID:    req.ParentID,
		Level:       level,
		Path:        path,
		PriceMin:    req.PriceMin,
		PriceMax:    req.PriceMax,
		TaxRate:     req.TaxRate,
		Template:    req.Template,
		Description: req.Description,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}

	if serviceType.Status == 0 {
		serviceType.Status = 1
	}
	if serviceType.TaxRate == 0 {
		serviceType.TaxRate = 6 // 默认6%税率
	}

	if err := database.DB.Create(&serviceType).Error; err != nil {
		return nil, err
	}

	return &serviceType, nil
}

// UpdateServiceType 更新服务类型
func (s *ServiceTypeService) UpdateServiceType(id uint, req *UpdateServiceTypeRequest) (*models.ServiceType, error) {
	var serviceType models.ServiceType
	if err := database.DB.First(&serviceType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务类型不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.PriceMin >= 0 {
		updates["price_min"] = req.PriceMin
	}
	if req.PriceMax >= 0 {
		updates["price_max"] = req.PriceMax
	}
	if req.TaxRate > 0 {
		updates["tax_rate"] = req.TaxRate
	}
	if req.Template != nil {
		updates["template"] = req.Template
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.SortOrder != 0 {
		updates["sort_order"] = req.SortOrder
	}

	if err := database.DB.Model(&serviceType).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &serviceType, nil
}

// DeleteServiceType 删除服务类型
func (s *ServiceTypeService) DeleteServiceType(id uint) error {
	var serviceType models.ServiceType
	if err := database.DB.First(&serviceType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("服务类型不存在")
		}
		return err
	}

	// 检查是否有子类型
	var childCount int64
	database.DB.Model(&models.ServiceType{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("该服务类型下有子类型，无法删除")
	}

	return database.DB.Delete(&serviceType).Error
}

// GetServiceType 获取服务类型详情
func (s *ServiceTypeService) GetServiceType(id uint) (*models.ServiceType, error) {
	var serviceType models.ServiceType
	if err := database.DB.Preload("Parent").Preload("Children").First(&serviceType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务类型不存在")
		}
		return nil, err
	}
	return &serviceType, nil
}

// GetServiceTypeByCode 根据代码获取服务类型
func (s *ServiceTypeService) GetServiceTypeByCode(code string) (*models.ServiceType, error) {
	var serviceType models.ServiceType
	if err := database.DB.Where("code = ?", code).First(&serviceType).Error; err != nil {
		return nil, err
	}
	return &serviceType, nil
}

// ListServiceTypes 获取服务类型列表
func (s *ServiceTypeService) ListServiceTypes(req *ListServiceTypesRequest) (*ListServiceTypesResponse, error) {
	var serviceTypes []models.ServiceType
	var total int64

	query := database.DB.Model(&models.ServiceType{})

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("sort_order ASC, created_at DESC").Find(&serviceTypes).Error; err != nil {
		return nil, err
	}

	return &ListServiceTypesResponse{
		List:  serviceTypes,
		Total: total,
	}, nil
}

// GetAllServiceTypes 获取所有服务类型（树形结构）
func (s *ServiceTypeService) GetAllServiceTypes() ([]models.ServiceType, error) {
	var serviceTypes []models.ServiceType
	if err := database.DB.Where("status = ? AND parent_id IS NULL", 1).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", 1).Order("sort_order ASC")
		}).
		Order("sort_order ASC").Find(&serviceTypes).Error; err != nil {
		return nil, err
	}
	return serviceTypes, nil
}

// GetServiceTypeTree 获取服务类型树
func (s *ServiceTypeService) GetServiceTypeTree() ([]map[string]interface{}, error) {
	var allTypes []models.ServiceType
	if err := database.DB.Where("status = ?", 1).Order("sort_order ASC").Find(&allTypes).Error; err != nil {
		return nil, err
	}

	// 构建树形结构
	nodeMap := make(map[uint]map[string]interface{})
	var roots []map[string]interface{}

	for _, t := range allTypes {
		node := map[string]interface{}{
			"id":          t.ID,
			"name":        t.Name,
			"code":        t.Code,
			"level":       t.Level,
			"price_min":   t.PriceMin,
			"price_max":   t.PriceMax,
			"tax_rate":    t.TaxRate,
			"children":    []map[string]interface{}{},
		}
		nodeMap[t.ID] = node

		if t.ParentID == nil {
			roots = append(roots, node)
		}
	}

	for _, t := range allTypes {
		if t.ParentID != nil {
			if parent, ok := nodeMap[*t.ParentID]; ok {
				children := parent["children"].([]map[string]interface{})
				children = append(children, nodeMap[t.ID])
				parent["children"] = children
			}
		}
	}

	return roots, nil
}
