import { get, post, put, del } from './request';
import type { 
  SystemConfig,
  OperationLog,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 系统设置 ====================

// 获取系统配置列表
export const getConfigList = (params?: PageParams & { group?: string }): Promise<PageResult<SystemConfig>> => {
  return get<PageResult<SystemConfig>>('/settings/configs', params as Record<string, unknown>);
};

// 获取配置详情
export const getConfigById = (id: number): Promise<SystemConfig> => {
  return get<SystemConfig>(`/settings/configs/${id}`);
};

// 根据key获取配置
export const getConfigByKey = (key: string): Promise<string> => {
  return get<string>(`/settings/configs/key/${key}`);
};

// 创建配置
export const createConfig = (data: { key: string; value: string; description?: string; group: string }): Promise<SystemConfig> => {
  return post<SystemConfig>('/settings/configs', data as Record<string, unknown>);
};

// 更新配置
export const updateConfig = (id: number, data: { value: string; description?: string }): Promise<SystemConfig> => {
  return put<SystemConfig>(`/settings/configs/${id}`, data as Record<string, unknown>);
};

// 删除配置
export const deleteConfig = (id: number): Promise<void> => {
  return del<void>(`/settings/configs/${id}`);
};

// 批量更新配置
export const batchUpdateConfigs = (configs: { key: string; value: string }[]): Promise<void> => {
  return post<void>('/settings/configs/batch-update', { configs });
};

// ==================== 操作日志 ====================

// 获取操作日志列表
export const getOperationLogList = (params?: PageParams & { 
  module?: string;
  operation?: string;
  operatorId?: number;
  status?: number;
  startDate?: string;
  endDate?: string;
}): Promise<PageResult<OperationLog>> => {
  return get<PageResult<OperationLog>>('/settings/logs', params as Record<string, unknown>);
};

// 获取操作日志详情
export const getOperationLogById = (id: number): Promise<OperationLog> => {
  return get<OperationLog>(`/settings/logs/${id}`);
};

// 清空日志
export const clearOperationLogs = (params?: { startDate?: string; endDate?: string }): Promise<void> => {
  return post<void>('/settings/logs/clear', params as Record<string, unknown>);
};
