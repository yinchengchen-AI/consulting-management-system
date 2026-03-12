import { get, post, put, del } from './request';
import type { 
  Customer, 
  CustomerFormData, 
  FollowUpRecord,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 客户管理 ====================

// 获取客户列表
export const getCustomerList = (params?: PageParams & { 
  keyword?: string; 
  status?: number; 
  industry?: string;
  source?: string;
  tag?: string;
}): Promise<PageResult<Customer>> => {
  return get<PageResult<Customer>>('/customers', params as Record<string, unknown>);
};

// 获取客户详情
export const getCustomerById = (id: number): Promise<Customer> => {
  return get<Customer>(`/customers/${id}`);
};

// 创建客户
export const createCustomer = (data: CustomerFormData): Promise<Customer> => {
  return post<Customer>('/customers', data as Record<string, unknown>);
};

// 更新客户
export const updateCustomer = (id: number, data: CustomerFormData): Promise<Customer> => {
  return put<Customer>(`/customers/${id}`, data as Record<string, unknown>);
};

// 删除客户
export const deleteCustomer = (id: number): Promise<void> => {
  return del<void>(`/customers/${id}`);
};

// 批量删除客户
export const batchDeleteCustomers = (ids: number[]): Promise<void> => {
  return post<void>('/customers/batch-delete', { ids });
};

// 获取所有客户（不分页）
export const getAllCustomers = (): Promise<Customer[]> => {
  return get<Customer[]>('/customers/all');
};

// 获取客户标签列表
export const getCustomerTags = (): Promise<string[]> => {
  return get<string[]>('/customers/tags');
};

// 获取行业列表
export const getIndustries = (): Promise<string[]> => {
  return get<string[]>('/customers/industries');
};

// ==================== 跟进记录 ====================

// 获取跟进记录列表
export const getFollowUpList = (params?: PageParams & { customerId?: number }): Promise<PageResult<FollowUpRecord>> => {
  return get<PageResult<FollowUpRecord>>('/follow-ups', params as Record<string, unknown>);
};

// 创建跟进记录
export const createFollowUp = (data: { customerId: number; content: string; type: string; nextFollowUpTime?: string }): Promise<FollowUpRecord> => {
  return post<FollowUpRecord>('/follow-ups', data as Record<string, unknown>);
};

// 更新跟进记录
export const updateFollowUp = (id: number, data: { content: string; type: string; nextFollowUpTime?: string }): Promise<FollowUpRecord> => {
  return put<FollowUpRecord>(`/follow-ups/${id}`, data as Record<string, unknown>);
};

// 删除跟进记录
export const deleteFollowUp = (id: number): Promise<void> => {
  return del<void>(`/follow-ups/${id}`);
};

// 获取客户的跟进记录
export const getCustomerFollowUps = (customerId: number): Promise<FollowUpRecord[]> => {
  return get<FollowUpRecord[]>(`/customers/${customerId}/follow-ups`);
};
