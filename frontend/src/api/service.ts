import { get, post, put, del } from './request';
import type { 
  ServiceType, 
  ServiceTypeFormData,
  ServiceOrder,
  ServiceOrderFormData,
  ServiceProgress,
  CommunicationRecord,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 服务类型管理 ====================

// 获取服务类型列表
export const getServiceTypeList = (params?: PageParams & { keyword?: string; status?: number }): Promise<PageResult<ServiceType>> => {
  return get<PageResult<ServiceType>>('/service-types', params as Record<string, unknown>);
};

// 获取所有服务类型（不分页）
export const getAllServiceTypes = (): Promise<ServiceType[]> => {
  return get<ServiceType[]>('/service-types/all');
};

// 获取服务类型详情
export const getServiceTypeById = (id: number): Promise<ServiceType> => {
  return get<ServiceType>(`/service-types/${id}`);
};

// 创建服务类型
export const createServiceType = (data: ServiceTypeFormData): Promise<ServiceType> => {
  return post<ServiceType>('/service-types', data as Record<string, unknown>);
};

// 更新服务类型
export const updateServiceType = (id: number, data: ServiceTypeFormData): Promise<ServiceType> => {
  return put<ServiceType>(`/service-types/${id}`, data as Record<string, unknown>);
};

// 删除服务类型
export const deleteServiceType = (id: number): Promise<void> => {
  return del<void>(`/service-types/${id}`);
};

// ==================== 服务订单管理 ====================

// 获取服务订单列表
export const getServiceOrderList = (params?: PageParams & { 
  keyword?: string; 
  status?: string;
  customerId?: number;
  serviceTypeId?: number;
  managerId?: number;
}): Promise<PageResult<ServiceOrder>> => {
  return get<PageResult<ServiceOrder>>('/service-orders', params as Record<string, unknown>);
};

// 获取服务订单详情
export const getServiceOrderById = (id: number): Promise<ServiceOrder> => {
  return get<ServiceOrder>(`/service-orders/${id}`);
};

// 创建服务订单
export const createServiceOrder = (data: ServiceOrderFormData): Promise<ServiceOrder> => {
  return post<ServiceOrder>('/service-orders', data as Record<string, unknown>);
};

// 更新服务订单
export const updateServiceOrder = (id: number, data: ServiceOrderFormData): Promise<ServiceOrder> => {
  return put<ServiceOrder>(`/service-orders/${id}`, data as Record<string, unknown>);
};

// 删除服务订单
export const deleteServiceOrder = (id: number): Promise<void> => {
  return del<void>(`/service-orders/${id}`);
};

// 更新服务订单状态
export const updateServiceOrderStatus = (id: number, status: string): Promise<void> => {
  return put<void>(`/service-orders/${id}/status`, { status });
};

// 更新服务订单进度
export const updateServiceOrderProgress = (id: number, progress: number): Promise<void> => {
  return put<void>(`/service-orders/${id}/progress`, { progress });
};

// ==================== 服务进度管理 ====================

// 获取服务进度列表
export const getServiceProgressList = (orderId: number): Promise<ServiceProgress[]> => {
  return get<ServiceProgress[]>(`/service-orders/${orderId}/progress`);
};

// 创建服务进度
export const createServiceProgress = (data: { orderId: number; title: string; content?: string; progress: number }): Promise<ServiceProgress> => {
  return post<ServiceProgress>('/service-progress', data as Record<string, unknown>);
};

// 删除服务进度
export const deleteServiceProgress = (id: number): Promise<void> => {
  return del<void>(`/service-progress/${id}`);
};

// ==================== 沟通记录管理 ====================

// 获取沟通记录列表
export const getCommunicationList = (orderId: number): Promise<CommunicationRecord[]> => {
  return get<CommunicationRecord[]>(`/service-orders/${orderId}/communications`);
};

// 创建沟通记录
export const createCommunication = (data: { orderId: number; content: string; communicationTime: string; participants?: string }): Promise<CommunicationRecord> => {
  return post<CommunicationRecord>('/communications', data as Record<string, unknown>);
};

// 删除沟通记录
export const deleteCommunication = (id: number): Promise<void> => {
  return del<void>(`/communications/${id}`);
};

// 简化的API函数名（兼容页面组件）
export const getServiceTypes = getServiceTypeList;
export const getServiceOrders = getServiceOrderList;
export const getServiceOrder = getServiceOrderById;
export const getCommunications = getCommunicationList;
