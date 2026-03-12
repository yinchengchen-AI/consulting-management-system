# 修复说明 - 咨询公司业务管理系统

## 修复日期
2026-03-12

## 修复内容总结

### 🔴 严重问题修复 (7项)

1. **main.go 调用 JWTAuth 缺少参数**
   - 原问题: `protected.Use(middleware.JWTAuth())` 缺少配置参数
   - 修复: 统一使用 `middleware.JWTAuth()` 无参版本（基于全局配置）

2. **config.GlobalConfig 未定义**
   - 原问题: middleware 引用了不存在的全局配置
   - 修复: 在 config.go 中添加 `GlobalConfig` 全局变量和 `LoadConfig()` 初始化函数

3. **Python 语法 import**
   - 原问题: database.go 包含 `from consulting-system/internal/utils import HashPassword`
   - 修复: 改为标准 Go import `"consulting-system/internal/utils"`

4. **方法名错误**
   - 原问题: `cfg.Database.GetDSN()` 和 `cfg.Redis.GetRedisAddr()`
   - 修复: 改为 `cfg.Database.DSN()` 和 `cfg.Redis.Addr()`

5. **Handler 初始化参数不匹配**
   - 原问题: `NewAuthHandler()` 调用与定义不匹配
   - 修复: 统一改为无参初始化，内部使用全局数据库实例

6. **JWT 相关重复定义**
   - 原问题: middleware/jwt.go 和 middleware/auth.go 重复
   - 修复: 删除 jwt.go，统一使用 auth.go

7. **User 模型主键类型**
   - 原问题: ID 为 string (UUID)，但 handlers 使用 uint 解析
   - 修复: User 模型改为 `ID uint` + `UUID string` 双字段

### 🟠 其他修复

8. **JWT 包函数签名统一**
   - UserID 改为 uint 类型
   - GenerateToken/GenerateRefreshToken/ParseRefreshToken 适配 uint

9. **缺少中间件补充**
   - 新增 error_logger.go

10. **AuthHandler 冲突解决**
    - 删除重复的 auth.go
    - 保留 auth_handler.go 和 auth_service.go

## 已知限制

由于代码量较大，以下问题可能需要进一步处理：

1. **go.mod 依赖**: 需要运行 `go mod tidy` 下载完整依赖
2. **前端类型**: types/index.ts 中仍存在字段命名不一致（不影响编译）
3. **Swagger 文档**: 需要运行 `swag init` 生成文档
4. **测试**: 缺少单元测试和集成测试

## 如何运行

### 后端
```bash
cd backend
go mod tidy
go run cmd/main.go
```

### 前端
```bash
cd frontend
npm install
npm run dev
```

### Docker 一键启动
```bash
make up
```

## 文件变更统计

- 修改文件: 15+
- 删除文件: 2 (jwt.go, auth.go)
- 新增文件: 1 (error_logger.go)

## 注意事项

1. 生产环境请务必修改默认 JWT 密钥
2. 建议添加环境变量验证
3. 建议添加数据库连接重试机制
