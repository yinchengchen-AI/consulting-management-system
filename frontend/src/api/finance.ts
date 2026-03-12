import { get, post, put, del } from './request';
import type { 
  Invoice, 
  InvoiceFormData,
  PaymentPlan,
  PaymentPlanFormData,
  Receipt,
  ReceiptFormData,
  Refund,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 开票管理 ====================

// 获取开票列表
export const getInvoiceList = (params?: PageParams & { 
  keyword?: string; 
  status?: string;
  customerId?: number;
  type?: string;
  startDate?: string;
  endDate?: string;
}): Promise<PageResult<Invoice>> => {
  return get<PageResult<Invoice>>('/invoices', params as Record<string, unknown>);
};

// 获取开票详情
export const getInvoiceById = (id: number): Promise<Invoice> => {
  return get<Invoice>(`/invoices/${id}`);
};

// 创建开票
export const createInvoice = (data: InvoiceFormData): Promise<Invoice> => {
  return post<Invoice>('/invoices', data as Record<string, unknown>);
};

// 更新开票
export const updateInvoice = (id: number, data: InvoiceFormData): Promise<Invoice> => {
  return put<Invoice>(`/invoices/${id}`, data as Record<string, unknown>);
};

// 删除开票
export const deleteInvoice = (id: number): Promise<void> => {
  return del<void>(`/invoices/${id}`);
};

// 提交审核
export const submitInvoice = (id: number): Promise<void> => {
  return put<void>(`/invoices/${id}/submit`);
};

// 审核开票
export const approveInvoice = (id: number, data: { approved: boolean; remark?: string }): Promise<void> => {
  return put<void>(`/invoices/${id}/approve`, data as Record<string, unknown>);
};

// 开具发票
export const issueInvoice = (id: number, data: { invoiceDate: string; remark?: string }): Promise<void> => {
  return put<void>(`/invoices/${id}/issue`, data as Record<string, unknown>);
};

// 取消开票
export const cancelInvoice = (id: number, reason: string): Promise<void> => {
  return put<void>(`/invoices/${id}/cancel`, { reason });
};

// ==================== 收款计划管理 ====================

// 获取收款计划列表
export const getPaymentPlanList = (params?: PageParams & { 
  keyword?: string; 
  status?: string;
  customerId?: number;
}): Promise<PageResult<PaymentPlan>> => {
  return get<PageResult<PaymentPlan>>('/payment-plans', params as Record<string, unknown>);
};

// 获取收款计划详情
export const getPaymentPlanById = (id: number): Promise<PaymentPlan> => {
  return get<PaymentPlan>(`/payment-plans/${id}`);
};

// 创建收款计划
export const createPaymentPlan = (data: PaymentPlanFormData): Promise<PaymentPlan> => {
  return post<PaymentPlan>('/payment-plans', data as Record<string, unknown>);
};

// 更新收款计划
export const updatePaymentPlan = (id: number, data: PaymentPlanFormData): Promise<PaymentPlan> => {
  return put<PaymentPlan>(`/payment-plans/${id}`, data as Record<string, unknown>);
};

// 删除收款计划
export const deletePaymentPlan = (id: number): Promise<void> => {
  return del<void>(`/payment-plans/${id}`);
};

// ==================== 收款记录管理 ====================

// 获取收款记录列表
export const getReceiptList = (params?: PageParams & { 
  keyword?: string; 
  planId?: number;
  customerId?: number;
  startDate?: string;
  endDate?: string;
}): Promise<PageResult<Receipt>> => {
  return get<PageResult<Receipt>>('/receipts', params as Record<string, unknown>);
};

// 获取收款记录详情
export const getReceiptById = (id: number): Promise<Receipt> => {
  return get<Receipt>(`/receipts/${id}`);
};

// 创建收款记录
export const createReceipt = (data: ReceiptFormData): Promise<Receipt> => {
  return post<Receipt>('/receipts', data as Record<string, unknown>);
};

// 更新收款记录
export const updateReceipt = (id: number, data: ReceiptFormData): Promise<Receipt> => {
  return put<Receipt>(`/receipts/${id}`, data as Record<string, unknown>);
};

// 删除收款记录
export const deleteReceipt = (id: number): Promise<void> => {
  return del<void>(`/receipts/${id}`);
};

// ==================== 退款管理 ====================

// 获取退款记录列表
export const getRefundList = (params?: PageParams & { receiptId?: number }): Promise<PageResult<Refund>> => {
  return get<PageResult<Refund>>('/refunds', params as Record<string, unknown>);
};

// 创建退款记录
export const createRefund = (data: { receiptId: number; amount: number; reason: string }): Promise<Refund> => {
  return post<Refund>('/refunds', data as Record<string, unknown>);
};

// 审核退款
export const approveRefund = (id: number, data: { approved: boolean; remark?: string }): Promise<void> => {
  return put<void>(`/refunds/${id}/approve`, data as Record<string, unknown>);
};

// ==================== 简化的API函数名（兼容页面组件）====================

export const getInvoices = getInvoiceList;
export const getInvoice = getInvoiceById;
export const auditInvoice = approveInvoice;

export const getPaymentPlans = getPaymentPlanList;
export const getPaymentPlan = getPaymentPlanById;

export const getReceipts = getReceiptList;
export const getReceipt = getReceiptById;
