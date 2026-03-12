import { get, post, put, del } from './request';
import type { 
  Contract, 
  ContractFormData,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 合同管理 ====================

// 获取合同列表
export const getContractList = (params?: PageParams & { 
  keyword?: string; 
  status?: string;
  customerId?: number;
  type?: string;
  startDate?: string;
  endDate?: string;
}): Promise<PageResult<Contract>> => {
  return get<PageResult<Contract>>('/contracts', params as Record<string, unknown>);
};

// 获取合同详情
export const getContractById = (id: number): Promise<Contract> => {
  return get<Contract>(`/contracts/${id}`);
};

// 创建合同
export const createContract = (data: ContractFormData): Promise<Contract> => {
  return post<Contract>('/contracts', data as Record<string, unknown>);
};

// 更新合同
export const updateContract = (id: number, data: ContractFormData): Promise<Contract> => {
  return put<Contract>(`/contracts/${id}`, data as Record<string, unknown>);
};

// 删除合同
export const deleteContract = (id: number): Promise<void> => {
  return del<void>(`/contracts/${id}`);
};

// 批量删除合同
export const batchDeleteContracts = (ids: number[]): Promise<void> => {
  return post<void>('/contracts/batch-delete', { ids });
};

// 更新合同状态
export const updateContractStatus = (id: number, status: string): Promise<void> => {
  return put<void>(`/contracts/${id}/status`, { status });
};

// 获取合同类型列表
export const getContractTypes = (): Promise<string[]> => {
  return get<string[]>('/contracts/types');
};
