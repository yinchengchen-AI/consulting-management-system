package services

import (
	"errors"
	"time"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"gorm.io/gorm"
)

// NoticeService 通知服务
type NoticeService struct{}

// NewNoticeService 创建通知服务实例
func NewNoticeService() *NoticeService {
	return &NoticeService{}
}

// CreateNoticeRequest 创建通知请求
type CreateNoticeRequest struct {
	Title       string              `json:"title" binding:"required,max=200"`
	Content     string              `json:"content" binding:"required"`
	Type        int                 `json:"type" binding:"omitempty,oneof=1 2 3"`
	TargetRoles models.StringArray  `json:"target_roles" binding:"omitempty"`
	IsTop       bool                `json:"is_top" binding:"omitempty"`
}

// UpdateNoticeRequest 更新通知请求
type UpdateNoticeRequest struct {
	Title       string              `json:"title" binding:"omitempty,max=200"`
	Content     string              `json:"content" binding:"omitempty"`
	Type        int                 `json:"type" binding:"omitempty,oneof=1 2 3"`
	TargetRoles models.StringArray  `json:"target_roles" binding:"omitempty"`
	IsTop       bool                `json:"is_top" binding:"omitempty"`
}

// ListNoticesRequest 通知列表请求
type ListNoticesRequest struct {
	Title  string `form:"title"`
	Type   int    `form:"type"`
	Page   int    `form:"page,default=1"`
	PageSize int  `form:"page_size,default=10"`
}

// ListNoticesResponse 通知列表响应
type ListNoticesResponse struct {
	List  []models.Notice `json:"list"`
	Total int64           `json:"total"`
}

// CreateNotice 创建通知
func (s *NoticeService) CreateNotice(req *CreateNoticeRequest, createdBy uint) (*models.Notice, error) {
	notice := models.Notice{
		Title:       req.Title,
		Content:     req.Content,
		Type:        req.Type,
		TargetRoles: req.TargetRoles,
		IsTop:       req.IsTop,
		CreatedBy:   createdBy,
	}

	if notice.Type == 0 {
		notice.Type = 1
	}

	if err := database.DB.Create(&notice).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("User").First(&notice, notice.ID)

	return &notice, nil
}

// UpdateNotice 更新通知
func (s *NoticeService) UpdateNotice(id uint, req *UpdateNoticeRequest) (*models.Notice, error) {
	var notice models.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("通知不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Type != 0 {
		updates["type"] = req.Type
	}
	if req.TargetRoles != nil {
		updates["target_roles"] = req.TargetRoles
	}
	updates["is_top"] = req.IsTop

	if err := database.DB.Model(&notice).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("User").First(&notice, id)

	return &notice, nil
}

// DeleteNotice 删除通知
func (s *NoticeService) DeleteNotice(id uint) error {
	var notice models.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return err
	}

	// 删除相关的阅读记录
	database.DB.Where("notice_id = ?", id).Delete(&models.NoticeRead{})

	return database.DB.Delete(&notice).Error
}

// GetNotice 获取通知详情
func (s *NoticeService) GetNotice(id uint, userID uint) (*models.Notice, error) {
	var notice models.Notice
	if err := database.DB.Preload("User").First(&notice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("通知不存在")
		}
		return nil, err
	}

	// 标记为已读
	if userID > 0 {
		s.MarkAsRead(id, userID)
	}

	return &notice, nil
}

// ListNotices 获取通知列表
func (s *NoticeService) ListNotices(req *ListNoticesRequest, userRoles []string) (*ListNoticesResponse, error) {
	var notices []models.Notice
	var total int64

	query := database.DB.Model(&models.Notice{})

	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Type != 0 {
		query = query.Where("type = ?", req.Type)
	}

	// 根据角色过滤
	if len(userRoles) > 0 && !contains(userRoles, "super_admin") {
		query = query.Where("target_roles @> ? OR target_roles = '[]'::jsonb OR target_roles IS NULL", 
			models.StringArray(userRoles))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("User").
		Offset(offset).Limit(req.PageSize).
		Order("is_top DESC, created_at DESC").
		Find(&notices).Error; err != nil {
		return nil, err
	}

	return &ListNoticesResponse{
		List:  notices,
		Total: total,
	}, nil
}

// MarkAsRead 标记通知为已读
func (s *NoticeService) MarkAsRead(noticeID, userID uint) error {
	// 检查是否已存在阅读记录
	var existing models.NoticeRead
	if err := database.DB.Where("notice_id = ? AND user_id = ?", noticeID, userID).First(&existing).Error; err == nil {
		return nil // 已存在，不需要重复记录
	}

	read := models.NoticeRead{
		NoticeID: noticeID,
		UserID:   userID,
		ReadAt:   time.Now(),
	}

	return database.DB.Create(&read).Error
}

// GetUnreadCount 获取未读通知数量
func (s *NoticeService) GetUnreadCount(userID uint, userRoles []string) (int64, error) {
	var count int64

	query := database.DB.Model(&models.Notice{}).
		Joins("LEFT JOIN notice_reads ON notices.id = notice_reads.notice_id AND notice_reads.user_id = ?", userID).
		Where("notice_reads.id IS NULL")

	// 根据角色过滤
	if len(userRoles) > 0 && !contains(userRoles, "super_admin") {
		query = query.Where("target_roles @> ? OR target_roles = '[]'::jsonb OR target_roles IS NULL",
			models.StringArray(userRoles))
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetReadRecords 获取阅读记录
func (s *NoticeService) GetReadRecords(noticeID uint) ([]models.NoticeRead, error) {
	var records []models.NoticeRead
	if err := database.DB.Where("notice_id = ?", noticeID).Preload("User").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// contains 检查字符串切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
