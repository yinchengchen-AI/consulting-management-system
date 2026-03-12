-- ============================================
-- 咨询公司业务管理系统 - 数据库初始化脚本
-- ============================================

-- 创建数据库用户（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'consulting_user') THEN
        CREATE USER consulting_user WITH PASSWORD 'consulting_pass';
    END IF;
END
$$;

-- 创建数据库（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'consulting_system') THEN
        CREATE DATABASE consulting_system OWNER consulting_user;
    END IF;
END
$$;

-- 连接到新创建的数据库
\c consulting_system;

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- 用于全文搜索

-- 授予用户权限
GRANT ALL PRIVILEGES ON DATABASE consulting_system TO consulting_user;

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建创建时间自动设置函数
CREATE OR REPLACE FUNCTION set_created_at_column()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.created_at IS NULL THEN
        NEW.created_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建软删除函数
CREATE OR REPLACE FUNCTION soft_delete()
RETURNS TRIGGER AS $$
BEGIN
    NEW.deleted_at = CURRENT_TIMESTAMP;
    NEW.is_deleted = true;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 设置搜索路径
ALTER DATABASE consulting_system SET search_path TO public;

-- 创建schema（用于组织表）
CREATE SCHEMA IF NOT EXISTS app;
GRANT ALL ON SCHEMA app TO consulting_user;
GRANT ALL ON SCHEMA public TO consulting_user;

-- 设置默认权限
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO consulting_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO consulting_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO consulting_user;

-- 初始化完成日志
DO $$
BEGIN
    RAISE NOTICE 'Database consulting_system initialized successfully!';
    RAISE NOTICE 'User consulting_user created with appropriate privileges.';
END
$$;
