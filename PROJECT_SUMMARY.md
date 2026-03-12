# 咨询公司业务管理系统 - 项目完成总结

## 项目概述

本项目是一个完整的**咨询公司业务管理系统**，采用 **Go + React + TypeScript** 技术栈，实现了文档中要求的所有11个核心模块。

## 技术栈

### 后端
- **Go 1.22** - 编程语言
- **Gin v1.9+** - Web框架
- **GORM v1.25+** - ORM框架
- **PostgreSQL 15** - 关系型数据库
- **Redis 7.x** - 缓存数据库
- **JWT v5** - 身份认证
- **Swagger** - API文档

### 前端
- **React 18** - UI框架
- **TypeScript 5.9+** - 类型系统
- **Vite 7.x** - 构建工具
- **Ant Design 5.29+** - UI组件库
- **Ant Design ProComponents 2.8+** - 高级组件
- **React Query 5.90+** - 数据获取
- **Zustand 5.x** - 状态管理
- **React Router 7.x** - 路由管理
- **ECharts 5.x** - 数据可视化
- **Axios 1.13+** - HTTP客户端

## 功能模块（11个）

### 1. 用户权限管理模块
- 用户CRUD管理
- 角色权限管理（RBAC）
- JWT认证登录/登出
- 密码安全管控

### 2. 客户公司管理模块
- 客户信息全生命周期管理
- 客户分类与标签
- 跟进记录管理
- 客户档案导出

### 3. 服务类型管理模块
- 服务类型配置
- 服务定价管理
- 服务模板管理
- 税率关联

### 4. 服务详细内容管理模块
- 服务单创建与管理
- 服务进度跟踪
- 沟通纪要记录
- 服务评价管理

### 5. 开票信息管理模块
- 开票申请流程
- 开票审核
- 发票归档管理
- 红字发票管理

### 6. 收款信息管理模块
- 收款计划管理
- 收款记录录入
- 退款管理
- 应收款统计

### 7. 统计分析模块
- 客户分析（数量、行业分布、高价值客户）
- 服务分析（类型占比、完成率）
- 财务分析（开票、收款、回款率）
- 数据可视化图表

### 8. 通知公告模块
- 通知发布
- 定向发送
- 阅读状态追踪

### 9. 文档管理模块
- 文档上传/下载
- 文档分类归档
- 访问权限控制

### 10. 合同管理模块
- 合同创建与管理
- 合同到期提醒
- 合同与服务、财务关联

### 11. 系统设置模块
- 系统参数配置
- 操作日志管理
- 数据备份与恢复

## 项目结构

```
consulting-system/
├── backend/                  # Go后端项目
│   ├── cmd/main.go          # 程序入口
│   ├── config/              # 配置文件
│   ├── internal/
│   │   ├── models/          # 数据库模型（11个模块）
│   │   ├── handlers/        # HTTP处理器
│   │   ├── services/        # 业务逻辑
│   │   ├── middleware/      # 中间件
│   │   ├── database/        # 数据库连接
│   │   └── utils/           # 工具函数
│   ├── pkg/                 # 公共包
│   ├── go.mod
│   └── Dockerfile
├── frontend/                 # React前端项目
│   ├── src/
│   │   ├── api/             # API调用（11个模块）
│   │   ├── components/      # 公共组件
│   │   ├── pages/           # 页面组件
│   │   ├── stores/          # 状态管理
│   │   ├── router/          # 路由配置
│   │   ├── types/           # 类型定义
│   │   └── utils/           # 工具函数
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── Dockerfile
├── docker/                   # Docker配置
│   ├── nginx/
│   └── postgres/
├── docker-compose.yml        # Docker编排
├── Makefile                  # 自动化脚本
└── README.md                 # 项目文档
```

## 文件统计

- **后端Go文件**: 49个
- **前端TS/TSX文件**: 48个
- **配置文件**: 4个
- **总计**: 100+ 个文件

## 快速开始

```bash
# 1. 进入项目目录
cd consulting-system

# 2. 启动所有服务
make up

# 3. 查看服务状态
make status

# 4. 查看日志
make logs
```

## 访问地址

- 前端应用: http://localhost
- 后端API: http://localhost:8080
- API文档: http://localhost:8080/swagger/index.html

## 核心特性

✅ 完整的RBAC权限控制
✅ JWT无状态认证
✅ RESTful API设计
✅ 响应式前端界面
✅ 数据可视化图表
✅ Docker容器化部署
✅ 数据库自动迁移
✅ 操作日志记录

## 开发团队

本项目严格按照文档要求开发，实现了咨询公司业务管理系统的全部功能。
