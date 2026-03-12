// ==================== 通用类型 ====================

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  size: number;
}

export interface PageParams {
  page?: number;
  pageSize?: number;
  size?: number;
  sort?: string;
  order?: 'asc' | 'desc';
  keyword?: string;
}

// ==================== 认证相关 ====================

export interface LoginParams {
  username: string;
  password: string;
  remember?: boolean;
}

export interface LoginResult {
  token: string;
  refreshToken: string;
  expiresIn: number;
  user: UserInfo;
}

export interface UserInfo {
  id: number;
  username: string;
  realName: string;
  real_name?: string;
  avatar?: string;
  email?: string;
  phone?: string;
  department?: string;
  position?: string;
  roles: string[];
  permissions: string[];
  status: number;
  lastLoginTime?: string;
}

// ==================== 用户管理 ====================

export interface User {
  id: number;
  username: string;
  realName: string;
  real_name?: string;
  email?: string;
  phone?: string;
  avatar?: string;
  departmentId?: number;
  department?: string;
  position?: string;
  status: number;
  roles: Role[];
  createTime?: string;
  created_at?: string;
  updateTime?: string;
  updated_at?: string;
  lastLoginTime?: string;
}

export interface UserFormData {
  id?: number;
  username: string;
  realName: string;
  password?: string;
  email?: string;
  phone?: string;
  departmentId?: number;
  position?: string;
  status: number;
  roleIds: number[];
}

export interface Role {
  id: number;
  name: string;
  code: string;
  description?: string;
  status: number;
  permissions: Permission[];
  createTime?: string;
  updateTime?: string;
}

export interface RoleFormData {
  id?: number;
  name: string;
  code: string;
  description?: string;
  status: number;
  permissionIds: number[];
}

export interface Permission {
  id: number;
  name: string;
  code: string;
  type: 'menu' | 'button' | 'api';
  parentId?: number;
  path?: string;
  icon?: string;
  sort: number;
  status: number;
  children?: Permission[];
}

// ==================== 客户管理 ====================

export interface Customer {
  id: number;
  name: string;
  shortName?: string;
  industry?: string;
  scale?: string;
  website?: string;
  address?: string;
  contactName?: string;
  contact_name?: string;
  contactPhone?: string;
  contact_phone?: string;
  contactEmail?: string;
  contact_email?: string;
  status: number;
  source?: string;
  tags: string[];
  remark?: string;
  followUpCount?: number;
  lastFollowUpTime?: string;
  createTime?: string;
  created_at?: string;
  updateTime?: string;
  updated_at?: string;
}

export interface CustomerFormData {
  id?: number;
  name: string;
  shortName?: string;
  industry?: string;
  scale?: string;
  website?: string;
  address?: string;
  contactName?: string;
  contactPhone?: string;
  contactEmail?: string;
  status: number;
  source?: string;
  tags: string[];
  remark?: string;
}

export interface FollowUpRecord {
  id: number;
  customerId: number;
  customerName?: string;
  content: string;
  type: string;
  followUpTime: string;
  nextFollowUpTime?: string;
  operatorId: number;
  operatorName?: string;
  createTime?: string;
  created_at?: string;
}

// ==================== 服务类型管理 ====================

export interface ServiceType {
  id: number;
  name: string;
  code: string;
  parent_id?: number;
  level?: number;
  path?: string;
  priceMin?: number;
  price_min?: number;
  priceMax?: number;
  price_max?: number;
  tax_rate?: number;
  TaxRate?: number;
  template?: Record<string, unknown>;
  description?: string;
  status: number;
  sortOrder?: number;
  sort_order?: number;
  createTime?: string;
  created_at?: string;
  updateTime?: string;
  updated_at?: string;
  parent?: ServiceType;
  children?: ServiceType[];
}

export interface ServiceTypeFormData {
  id?: number;
  name: string;
  code: string;
  description?: string;
  price?: number;
  unit?: string;
  duration?: number;
  status: number;
  sort: number;
}

// ==================== 服务订单管理 ====================

export interface ServiceOrder {
  id: number;
  code?: string;
  orderNo?: string;
  customer_id?: number;
  CustomerID?: number;
  customerId?: number;
  customer?: Customer;
  service_type_id?: number;
  ServiceTypeID?: number;
  serviceTypeId?: number;
  service_type?: ServiceType;
  ServiceType?: ServiceType;
  name?: string;
  title?: string;
  description?: string;
  start_date?: string;
  StartDate?: string;
  startDate?: string;
  end_date?: string;
  EndDate?: string;
  endDate?: string;
  amount?: number;
  Amount?: number;
  status: number;
  progress?: number;
  Progress?: number;
  participants?: string[];
  Participants?: string[];
  managerId?: number;
  managerName?: string;
  created_by?: number;
  CreatedBy?: number;
  createTime?: string;
  created_at?: string;
  updateTime?: string;
  updated_at?: string;
  communications?: Communication[];
}

export type ServiceOrderStatus = 
  | 'draft'
  | 'pending'
  | 'processing'
  | 'completed'
  | 'cancelled';

export interface ServiceOrderFormData {
  id?: number;
  customer_id?: number;
  service_type_id?: number;
  code?: string;
  name?: string;
  title?: string;
  description?: string;
  start_date?: string;
  end_date?: string;
  amount?: number;
  status?: number;
  progress?: number;
  participants?: string[];
}

export interface ServiceProgress {
  id: number;
  orderId: number;
  title: string;
  content?: string;
  progress: number;
  operatorId: number;
  operatorName?: string;
  createTime?: string;
}

export interface Communication {
  id: number;
  service_id?: number;
  ServiceID?: number;
  orderId?: number;
  content: string;
  communicationTime?: string;
  participants?: string;
  created_by?: number;
  CreatedBy?: number;
  created_at?: string;
  createTime?: string;
  user?: User;
}

export interface CommunicationRecord {
  id: number;
  orderId: number;
  content: string;
  communicationTime: string;
  participants?: string;
  operatorId: number;
  operatorName?: string;
  createTime?: string;
}

// ==================== 开票管理 ====================

export interface Invoice {
  id: number;
  customer_id?: number;
  CustomerID?: number;
  customerId?: number;
  customer?: Customer;
  service_id?: number;
  ServiceID?: number;
  serviceId?: number;
  invoice_type?: number;
  InvoiceType?: number;
  invoiceType?: number;
  amount?: number;
  Amount?: number;
  tax_rate?: number;
  TaxRate?: number;
  tax_amount?: number;
  TaxAmount?: number;
  total_amount?: number;
  TotalAmount?: number;
  invoice_no?: string;
  InvoiceNo?: string;
  invoiceNo?: string;
  invoice_code?: string;
  InvoiceCode?: string;
  invoiceCode?: string;
  status: number;
  invoice_info?: Record<string, unknown>;
  InvoiceInfo?: Record<string, unknown>;
  invoice_date?: string;
  InvoiceDate?: string;
  invoiceDate?: string;
  created_by?: number;
  CreatedBy?: number;
  created_at?: string;
  updated_at?: string;
}

export interface InvoiceFormData {
  id?: number;
  customerId?: number;
  serviceId?: number;
  invoiceType?: number;
  amount?: number;
  taxRate?: number;
  invoiceNo?: string;
  invoiceCode?: string;
  invoiceDate?: string;
  status?: number;
}

// ==================== 收款管理 ====================

export interface PaymentPlan {
  id: number;
  customer_id?: number;
  CustomerID?: number;
  customerId?: number;
  customer?: Customer;
  service_id?: number;
  ServiceID?: number;
  serviceId?: number;
  invoice_id?: number;
  InvoiceID?: number;
  invoiceId?: number;
  amount?: number;
  Amount?: number;
  planned_date?: string;
  PlannedDate?: string;
  plannedDate?: string;
  status?: number;
  Status?: number;
  remark?: string;
  Remark?: string;
  created_by?: number;
  CreatedBy?: number;
  created_at?: string;
  receipts?: Receipt[];
}

export interface PaymentPlanFormData {
  id?: number;
  customerId?: number;
  serviceId?: number;
  invoiceId?: number;
  amount?: number;
  plannedDate?: string;
  status?: number;
  remark?: string;
}

export interface Receipt {
  id: number;
  plan_id?: number;
  PlanID?: number;
  planId?: number;
  plan?: PaymentPlan;
  amount?: number;
  Amount?: number;
  received_date?: string;
  ReceivedDate?: string;
  receivedDate?: string;
  payment_method?: number;
  PaymentMethod?: number;
  paymentMethod?: number;
  account?: string;
  Account?: string;
  remark?: string;
  Remark?: string;
  created_by?: number;
  CreatedBy?: number;
  created_at?: string;
  refunds?: Refund[];
}

export interface ReceiptFormData {
  id?: number;
  planId?: number;
  amount?: number;
  receivedDate?: string;
  paymentMethod?: number;
  account?: string;
  remark?: string;
}

export interface Refund {
  id: number;
  receipt_id?: number;
  ReceiptID?: number;
  receiptId?: number;
  receipt?: Receipt;
  amount?: number;
  Amount?: number;
  reason?: string;
  Reason?: string;
  status?: number;
  Status?: number;
  created_by?: number;
  CreatedBy?: number;
  created_at?: string;
}

// ==================== 合同管理 ====================

export interface Contract {
  id: number;
  code?: string;
  name?: string;
  customer_id?: number;
  customerId?: number;
  customer?: Customer;
  service_id?: number;
  serviceId?: number;
  amount?: number;
  sign_date?: string;
  signDate?: string;
  expire_date?: string;
  expireDate?: string;
  payment_terms?: string;
  PaymentTerms?: string;
  status?: number;
  Status?: number;
  file_url?: string;
  FileURL?: string;
  fileUrl?: string;
  remark?: string;
  Remark?: string;
  created_by?: number;
  CreatedBy?: number;
  created_at?: string;
  updated_at?: string;
}

export interface ContractFormData {
  id?: number;
  code?: string;
  name?: string;
  customerId?: number;
  serviceId?: number;
  amount?: number;
  signDate?: string;
  expireDate?: string;
  paymentTerms?: string;
  status?: number;
  fileUrl?: string;
  remark?: string;
}

// ==================== 统计分析 ====================

export interface CustomerStats {
  total: number;
  newThisMonth: number;
  industryDistribution: { name: string; value: number }[];
  topCustomers: { name: string; amount: number }[];
}

export interface ServiceStats {
  total: number;
  inProgress: number;
  completed: number;
  typeDistribution: { name: string; value: number }[];
}

export interface FinanceStats {
  totalInvoice: number;
  totalReceipt: number;
  collectionRate: number;
  overdueAmount: number;
  monthlyTrend: { month: string; invoice: number; receipt: number }[];
}

// ==================== 通知公告 ====================

export interface Notice {
  id: number;
  title: string;
  content: string;
  type: number;
  targetRoles?: string[];
  isTop?: boolean;
  createdBy?: number;
  createdAt?: string;
}

// ==================== 文档管理 ====================

export interface Document {
  id: number;
  name: string;
  type?: string;
  fileUrl?: string;
  size?: number;
  mimeType?: string;
  relatedType?: string;
  relatedId?: number;
  accessPermissions?: Record<string, unknown>;
  description?: string;
  createdBy?: number;
  createdAt?: string;
}
