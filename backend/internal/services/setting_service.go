package services

import (
	"errors"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"gorm.io/gorm"
)

// SettingService 设置服务
type SettingService struct{}

// NewSettingService 创建设置服务实例
func NewSettingService() *SettingService {
	return &SettingService{}
}

// CreateConfigRequest 创建配置请求
type CreateConfigRequest struct {
	Key         string `json:"key" binding:"required,max=100"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	Value       string `json:"value" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// ListLogsRequest 日志列表请求
type ListLogsRequest struct {
	UserID     uint   `form:"user_id"`
	Action     string `form:"action"`
	TargetType string `form:"target_type"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
}

// ListLogsResponse 日志列表响应
type ListLogsResponse struct {
	List  []models.OperationLog `json:"list"`
	Total int64                 `json:"total"`
}

// CreateConfig 创建配置
func (s *SettingService) CreateConfig(req *CreateConfigRequest) (*models.SystemConfig, error) {
	// 检查key是否已存在
	var existing models.SystemConfig
	if err := database.DB.Where("key = ?", req.Key).First(&existing).Error; err == nil {
		return nil, errors.New("配置项已存在")
	}

	config := models.SystemConfig{
		Key:         req.Key,
		Value:       req.Value,
		Description: req.Description,
	}

	if err := database.DB.Create(&config).Error; err != nil {
		return nil, err
	}

	return &config, nil
}

// UpdateConfig 更新配置
func (s *SettingService) UpdateConfig(key string, req *UpdateConfigRequest) (*models.SystemConfig, error) {
	var config models.SystemConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("配置项不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Value != "" {
		updates["value"] = req.Value
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := database.DB.Model(&config).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &config, nil
}

// DeleteConfig 删除配置
func (s *SettingService) DeleteConfig(key string) error {
	var config models.SystemConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("配置项不存在")
		}
		return err
	}

	return database.DB.Delete(&config).Error
}

// GetConfig 获取配置
func (s *SettingService) GetConfig(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("配置项不存在")
		}
		return nil, err
	}
	return &config, nil
}

// GetConfigValue 获取配置值
func (s *SettingService) GetConfigValue(key string, defaultValue string) string {
	config, err := s.GetConfig(key)
	if err != nil {
		return defaultValue
	}
	return config.Value
}

// ListConfigs 获取所有配置
func (s *SettingService) ListConfigs() ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetConfigsMap 获取配置Map
func (s *SettingService) GetConfigsMap() (map[string]string, error) {
	configs, err := s.ListConfigs()
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	return configMap, nil
}

// ListLogs 获取操作日志列表
func (s *SettingService) ListLogs(req *ListLogsRequest) (*ListLogsResponse, error) {
	var logs []models.OperationLog
	var total int64

	query := database.DB.Model(&models.OperationLog{})

	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.Action != "" {
		query = query.Where("action LIKE ?", "%"+req.Action+"%")
	}
	if req.TargetType != "" {
		query = query.Where("target_type = ?", req.TargetType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("User").Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}

	return &ListLogsResponse{
		List:  logs,
		Total: total,
	}, nil
}

// GetLog 获取日志详情
func (s *SettingService) GetLog(id uint) (*models.OperationLog, error) {
	var log models.OperationLog
	if err := database.DB.Preload("User").First(&log, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("日志不存在")
		}
		return nil, err
	}
	return &log, nil
}

// DeleteLogs 删除日志
func (s *SettingService) DeleteLogs(beforeDate string) error {
	return database.DB.Where("created_at < ?", beforeDate).Delete(&models.OperationLog{}).Error
}

// CreateOperationLog 创建操作日志
func (s *SettingService) CreateOperationLog(userID uint, action, targetType string, targetID uint, details models.JSON, ip, userAgent string) error {
	log := models.OperationLog{
		UserID:     userID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Details:    details,
		IP:         ip,
		UserAgent:  userAgent,
	}

	return database.DB.Create(&log).Error
}
