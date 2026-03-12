package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

// GenerateCode 生成指定格式的编码
func GenerateCode(prefix string, seq int) string {
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("%s%s%06d", prefix, timestamp, seq)
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// IsValidPhone 验证手机号格式（中国）
func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(phone)
}

// SanitizeString 清理字符串
func SanitizeString(s string) string {
	// 去除前后空格
	s = strings.TrimSpace(s)
	// 去除多余空格
	spaceRegex := regexp.MustCompile(`\s+`)
	s = spaceRegex.ReplaceAllString(s, " ")
	return s
}

// TruncateString 截断字符串
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// ContainsString 检查字符串切片是否包含指定字符串
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates 去除字符串切片中的重复项
func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// FormatTime 格式化时间
func FormatTime(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string, layout string) (time.Time, error) {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return time.Parse(layout, timeStr)
}

// GetStartOfDay 获取当天的开始时间
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取当天的结束时间
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// GetStartOfMonth 获取当月的开始时间
func GetStartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfMonth 获取当月的结束时间
func GetEndOfMonth(t time.Time) time.Time {
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	lastDay := firstDay.AddDate(0, 1, -1)
	return time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 999999999, t.Location())
}

// CalculateAge 计算年龄
func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

// FormatMoney 格式化金额
func FormatMoney(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// Round 四舍五入
func Round(value float64, decimals int) float64 {
	format := fmt.Sprintf("%%.%df", decimals)
	result, _ := fmt.Sscanf(fmt.Sprintf(format, value), "%f", &value)
	if result != 1 {
		return value
	}
	return value
}

// IsEmptyString 检查字符串是否为空
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

// DefaultString 返回默认值
func DefaultString(value, defaultValue string) string {
	if IsEmptyString(value) {
		return defaultValue
	}
	return value
}
