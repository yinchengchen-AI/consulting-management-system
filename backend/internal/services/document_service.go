package services

import (
	"errors"
	"path/filepath"
	"strings"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentService 文档服务
type DocumentService struct{}

// NewDocumentService 创建文档服务实例
func NewDocumentService() *DocumentService {
	return &DocumentService{}
}

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	Name              string      `json:"name" binding:"required,max=255"`
	Type              string      `json:"type" binding:"omitempty,max=50"`
	FileURL           string      `json:"file_url" binding:"required"`
	Size              int64       `json:"size" binding:"omitempty,min=0"`
	MimeType          string      `json:"mime_type" binding:"omitempty,max=100"`
	RelatedType       string      `json:"related_type" binding:"omitempty,max=50"`
	RelatedID         uint        `json:"related_id" binding:"omitempty"`
	AccessPermissions models.JSON `json:"access_permissions" binding:"omitempty"`
	Description       string      `json:"description" binding:"omitempty"`
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Name              string      `json:"name" binding:"omitempty,max=255"`
	Type              string      `json:"type" binding:"omitempty,max=50"`
	AccessPermissions models.JSON `json:"access_permissions" binding:"omitempty"`
	Description       string      `json:"description" binding:"omitempty"`
}

// ListDocumentsRequest 文档列表请求
type ListDocumentsRequest struct {
	Name        string `form:"name"`
	Type        string `form:"type"`
	RelatedType string `form:"related_type"`
	RelatedID   uint   `form:"related_id"`
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
}

// ListDocumentsResponse 文档列表响应
type ListDocumentsResponse struct {
	List  []models.Document `json:"list"`
	Total int64             `json:"total"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// GenerateFileName 生成文件名
func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return uuid.New().String() + ext
}

// GetFileType 根据文件名获取文件类型
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "word"
	case ".xls", ".xlsx":
		return "excel"
	case ".ppt", ".pptx":
		return "ppt"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return "image"
	case ".mp4", ".avi", ".mov", ".wmv":
		return "video"
	case ".mp3", ".wav", ".flac":
		return "audio"
	case ".zip", ".rar", ".7z":
		return "archive"
	default:
		return "other"
	}
}

// CreateDocument 创建文档记录
func (s *DocumentService) CreateDocument(req *CreateDocumentRequest, createdBy uint) (*models.Document, error) {
	document := models.Document{
		Name:              req.Name,
		Type:              req.Type,
		FileURL:           req.FileURL,
		Size:              req.Size,
		MimeType:          req.MimeType,
		RelatedType:       req.RelatedType,
		RelatedID:         req.RelatedID,
		AccessPermissions: req.AccessPermissions,
		Description:       req.Description,
		CreatedBy:         createdBy,
	}

	if document.Type == "" {
		document.Type = GetFileType(req.Name)
	}

	if err := database.DB.Create(&document).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("User").First(&document, document.ID)

	return &document, nil
}

// UpdateDocument 更新文档
func (s *DocumentService) UpdateDocument(id uint, req *UpdateDocumentRequest) (*models.Document, error) {
	var document models.Document
	if err := database.DB.First(&document, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文档不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.AccessPermissions != nil {
		updates["access_permissions"] = req.AccessPermissions
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := database.DB.Model(&document).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("User").First(&document, id)

	return &document, nil
}

// DeleteDocument 删除文档
func (s *DocumentService) DeleteDocument(id uint) error {
	var document models.Document
	if err := database.DB.First(&document, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文档不存在")
		}
		return err
	}

	return database.DB.Delete(&document).Error
}

// GetDocument 获取文档详情
func (s *DocumentService) GetDocument(id uint) (*models.Document, error) {
	var document models.Document
	if err := database.DB.Preload("User").First(&document, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文档不存在")
		}
		return nil, err
	}
	return &document, nil
}

// ListDocuments 获取文档列表
func (s *DocumentService) ListDocuments(req *ListDocumentsRequest) (*ListDocumentsResponse, error) {
	var documents []models.Document
	var total int64

	query := database.DB.Model(&models.Document{})

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.RelatedType != "" {
		query = query.Where("related_type = ?", req.RelatedType)
	}
	if req.RelatedID != 0 {
		query = query.Where("related_id = ?", req.RelatedID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("User").Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, err
	}

	return &ListDocumentsResponse{
		List:  documents,
		Total: total,
	}, nil
}

// GetDocumentsByRelated 根据关联获取文档
func (s *DocumentService) GetDocumentsByRelated(relatedType string, relatedID uint) ([]models.Document, error) {
	var documents []models.Document
	if err := database.DB.Where("related_type = ? AND related_id = ?", relatedType, relatedID).
		Preload("User").Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

// CheckPermission 检查用户是否有权限访问文档
func (s *DocumentService) CheckPermission(documentID uint, userID uint, userRoles []string) bool {
	var document models.Document
	if err := database.DB.First(&document, documentID).Error; err != nil {
		return false
	}

	// 文档创建者有权限
	if document.CreatedBy == userID {
		return true
	}

	// 检查权限配置
	if document.AccessPermissions != nil {
		// 检查用户ID
		if userIDs, ok := document.AccessPermissions["user_ids"].([]interface{}); ok {
			for _, id := range userIDs {
				if uint(id.(float64)) == userID {
					return true
				}
			}
		}

		// 检查角色
		if roles, ok := document.AccessPermissions["roles"].([]interface{}); ok {
			for _, role := range roles {
				for _, userRole := range userRoles {
					if role.(string) == userRole {
						return true
					}
				}
			}
		}
	}

	return false
}
