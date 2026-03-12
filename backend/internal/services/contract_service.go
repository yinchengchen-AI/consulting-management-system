package services

import (
	"errors"
	"time"

	"consulting-system/internal/database"
	"consulting-system/internal/models"
	"consulting-system/internal/utils"

	"gorm.io/gorm"
)

// ContractService 合同服务
type ContractService struct{}

// NewContractService 创建合同服务实例
func NewContractService() *ContractService {
	return &ContractService{}
}

// CreateContractRequest 创建合同请求
type CreateContractRequest struct {
	CustomerID    uint       `json:"customer_id" binding:"required"`
	ServiceID     *uint      `json:"service_id" binding:"omitempty"`
	Name          string     `json:"name" binding:"required,max=200"`
	Amount        float64    `json:"amount" binding:"omitempty,min=0"`
	SignDate      *time.Time `json:"sign_date" binding:"omitempty"`
	ExpireDate    *time.Time `json:"expire_date" binding:"omitempty"`
	PaymentTerms  string     `json:"payment_terms" binding:"omitempty"`
	FileURL       string     `json:"file_url" binding:"omitempty"`
	Remark        string     `json:"remark" binding:"omitempty"`
}

// UpdateContractRequest 更新合同请求
type UpdateContractRequest struct {
	Name         string     `json:"name" binding:"omitempty,max=200"`
	Amount       float64    `json:"amount" binding:"omitempty,min=0"`
	SignDate     *time.Time `json:"sign_date" binding:"omitempty"`
	ExpireDate   *time.Time `json:"expire_date" binding:"omitempty"`
	PaymentTerms string     `json:"payment_terms" binding:"omitempty"`
	Status       int        `json:"status" binding:"omitempty,oneof=1 2 3 4 5"`
	FileURL      string     `json:"file_url" binding:"omitempty"`
	Remark       string     `json:"remark" binding:"omitempty"`
}

// ListContractsRequest 合同列表请求
type ListContractsRequest struct {
	CustomerID uint   `form:"customer_id"`
	Name       string `form:"name"`
	Status     int    `form:"status"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
}

// ListContractsResponse 合同列表响应
type ListContractsResponse struct {
	List  []models.Contract `json:"list"`
	Total int64             `json:"total"`
}

// generateContractCode 生成合同编号
func generateContractCode() string {
	seq := int(time.Now().UnixNano() % 1000000)
	return utils.GenerateCode("CTR", seq)
}

// CreateContract 创建合同
func (s *ContractService) CreateContract(req *CreateContractRequest, createdBy uint) (*models.Contract, error) {
	// 检查客户是否存在
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		return nil, errors.New("客户不存在")
	}

	// 检查服务订单是否存在
	if req.ServiceID != nil {
		var service models.ServiceOrder
		if err := database.DB.First(&service, *req.ServiceID).Error; err != nil {
			return nil, errors.New("服务订单不存在")
		}
	}

	contract := models.Contract{
		CustomerID:   req.CustomerID,
		ServiceID:    req.ServiceID,
		Code:         generateContractCode(),
		Name:         req.Name,
		Amount:       req.Amount,
		SignDate:     req.SignDate,
		ExpireDate:   req.ExpireDate,
		PaymentTerms: req.PaymentTerms,
		Status:       1, // 草稿
		FileURL:      req.FileURL,
		Remark:       req.Remark,
		CreatedBy:    createdBy,
	}

	if err := database.DB.Create(&contract).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&contract, contract.ID)

	return &contract, nil
}

// UpdateContract 更新合同
func (s *ContractService) UpdateContract(id uint, req *UpdateContractRequest) (*models.Contract, error) {
	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合同不存在")
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
	if req.SignDate != nil {
		updates["sign_date"] = req.SignDate
	}
	if req.ExpireDate != nil {
		updates["expire_date"] = req.ExpireDate
	}
	if req.PaymentTerms != "" {
		updates["payment_terms"] = req.PaymentTerms
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := database.DB.Model(&contract).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&contract, id)

	return &contract, nil
}

// DeleteContract 删除合同
func (s *ContractService) DeleteContract(id uint) error {
	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("合同不存在")
		}
		return err
	}

	return database.DB.Delete(&contract).Error
}

// GetContract 获取合同详情
func (s *ContractService) GetContract(id uint) (*models.Contract, error) {
	var contract models.Contract
	if err := database.DB.Preload("Customer").First(&contract, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合同不存在")
		}
		return nil, err
	}
	return &contract, nil
}

// ListContracts 获取合同列表
func (s *ContractService) ListContracts(req *ListContractsRequest) (*ListContractsResponse, error) {
	var contracts []models.Contract
	var total int64

	query := database.DB.Model(&models.Contract{})

	if req.CustomerID != 0 {
		query = query.Where("customer_id = ?", req.CustomerID)
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
	if err := query.Preload("Customer").Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&contracts).Error; err != nil {
		return nil, err
	}

	return &ListContractsResponse{
		List:  contracts,
		Total: total,
	}, nil
}

// GetExpiringContracts 获取即将到期的合同
func (s *ContractService) GetExpiringContracts(days int) ([]models.Contract, error) {
	var contracts []models.Contract

	expireDate := time.Now().AddDate(0, 0, days)

	if err := database.DB.Where("expire_date <= ? AND expire_date > ? AND status IN (?, ?)",
		expireDate, time.Now(), 2, 3).
		Preload("Customer").
		Order("expire_date ASC").
		Find(&contracts).Error; err != nil {
		return nil, err
	}

	return contracts, nil
}

// SignContract 签署合同
func (s *ContractService) SignContract(id uint) (*models.Contract, error) {
	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合同不存在")
		}
		return nil, err
	}

	if contract.Status != 1 {
		return nil, errors.New("只有草稿状态的合同可以签署")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":    2, // 已签署
		"sign_date": &now,
	}

	if err := database.DB.Model(&contract).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&contract, id)

	return &contract, nil
}

// TerminateContract 终止合同
func (s *ContractService) TerminateContract(id uint) (*models.Contract, error) {
	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合同不存在")
		}
		return nil, err
	}

	if contract.Status != 2 && contract.Status != 3 {
		return nil, errors.New("只有已签署或履行中的合同可以终止")
	}

	if err := database.DB.Model(&contract).Update("status", 5).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&contract, id)

	return &contract, nil
}

// GetContractStats 获取合同统计
func (s *ContractService) GetContractStats() (map[string]interface{}, error) {
	var total int64
	var byStatus []struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}
	var totalAmount float64
	var expiringCount int64

	database.DB.Model(&models.Contract{}).Count(&total)
	database.DB.Model(&models.Contract{}).Select("status, COUNT(*) as count").Group("status").Scan(&byStatus)
	database.DB.Model(&models.Contract{}).Select("COALESCE(SUM(amount), 0)").Pluck("COALESCE(SUM(amount), 0)", &totalAmount)
	database.DB.Model(&models.Contract{}).Where("expire_date <= ? AND expire_date > ? AND status IN (?, ?)",
		time.Now().AddDate(0, 0, 30), time.Now(), 2, 3).Count(&expiringCount)

	return map[string]interface{}{
		"total":          total,
		"by_status":      byStatus,
		"total_amount":   totalAmount,
		"expiring_count": expiringCount,
	}, nil
}
