import { get } from './request';
import type { 
  CustomerStats,
  ServiceStats,
  FinanceStats,
} from '@/types';

// ==================== 客户统计 ====================

// 获取客户统计数据
export const getCustomerStats = (params?: { startDate?: string; endDate?: string }): Promise<CustomerStats> => {
  return get<CustomerStats>('/statistics/customers', params as Record<string, unknown>);
};

// ==================== 服务统计 ====================

// 获取服务统计数据
export const getServiceStats = (params?: { startDate?: string; endDate?: string }): Promise<ServiceStats> => {
  return get<ServiceStats>('/statistics/services', params as Record<string, unknown>);
};

// ==================== 财务统计 ====================

// 获取财务统计数据
export const getFinanceStats = (params?: { startDate?: string; endDate?: string }): Promise<FinanceStats> => {
  return get<FinanceStats>('/statistics/finance', params as Record<string, unknown>);
};
