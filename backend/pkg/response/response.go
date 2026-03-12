package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseCode 响应状态码
type ResponseCode int

const (
	CodeSuccess       ResponseCode = 200
	CodeCreated       ResponseCode = 201
	CodeBadRequest    ResponseCode = 400
	CodeUnauthorized  ResponseCode = 401
	CodeForbidden     ResponseCode = 403
	CodeNotFound      ResponseCode = 404
	CodeConflict      ResponseCode = 409
	CodeValidationErr ResponseCode = 422
	CodeInternalErr   ResponseCode = 500
)

// ResponseMessage 响应消息
var ResponseMessage = map[ResponseCode]string{
	CodeSuccess:       "操作成功",
	CodeCreated:       "创建成功",
	CodeBadRequest:    "请求参数错误",
	CodeUnauthorized:  "未授权，请先登录",
	CodeForbidden:     "无权限访问",
	CodeNotFound:      "资源不存在",
	CodeConflict:      "资源冲突",
	CodeValidationErr: "数据验证失败",
	CodeInternalErr:   "服务器内部错误",
}

// Response 统一响应结构
type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
}

// Meta 分页元数据
type Meta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: ResponseMessage[CodeSuccess],
		Data:    data,
	})
}

// SuccessWithMeta 带分页的成功响应
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: ResponseMessage[CodeSuccess],
		Data:    data,
		Meta:    meta,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    CodeCreated,
		Message: ResponseMessage[CodeCreated],
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code ResponseCode, message string) {
	if message == "" {
		message = ResponseMessage[code]
	}
	c.JSON(getHTTPStatus(code), Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code ResponseCode, message string, data interface{}) {
	if message == "" {
		message = ResponseMessage[code]
	}
	c.JSON(getHTTPStatus(code), Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, CodeBadRequest, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = ResponseMessage[CodeUnauthorized]
	}
	Error(c, CodeUnauthorized, message)
}

// Forbidden 无权限
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = ResponseMessage[CodeForbidden]
	}
	Error(c, CodeForbidden, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = ResponseMessage[CodeNotFound]
	}
	Error(c, CodeNotFound, message)
}

// ValidationError 验证错误
func ValidationError(c *gin.Context, message string) {
	Error(c, CodeValidationErr, message)
}

// InternalError 服务器内部错误
func InternalError(c *gin.Context, message string) {
	Error(c, CodeInternalErr, message)
}

// getHTTPStatus 获取HTTP状态码
func getHTTPStatus(code ResponseCode) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeCreated:
		return http.StatusCreated
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeValidationErr:
		return http.StatusUnprocessableEntity
	case CodeInternalErr:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Pagination 分页参数
type Pagination struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *Pagination) GetLimit() int {
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// BuildMeta 构建分页元数据
func BuildMeta(total int64, page, pageSize int) *Meta {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &Meta{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
