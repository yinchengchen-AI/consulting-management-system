import { get, post, put, del } from './request';
import type { 
  User, 
  UserFormData, 
  Role, 
  RoleFormData, 
  Permission,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 用户管理 ====================

// 获取用户列表
export const getUserList = (params?: PageParams & { keyword?: string; status?: number; departmentId?: number }): Promise<PageResult<User>> => {
  return get<PageResult<User>>('/users', params as Record<string, unknown>);
};

// 获取用户详情
export const getUserById = (id: number): Promise<User> => {
  return get<User>(`/users/${id}`);
};

// 创建用户
export const createUser = (data: UserFormData): Promise<User> => {
  return post<User>('/users', data as Record<string, unknown>);
};

// 更新用户
export const updateUser = (id: number, data: UserFormData): Promise<User> => {
  return put<User>(`/users/${id}`, data as Record<string, unknown>);
};

// 删除用户
export const deleteUser = (id: number): Promise<void> => {
  return del<void>(`/users/${id}`);
};

// 批量删除用户
export const batchDeleteUsers = (ids: number[]): Promise<void> => {
  return post<void>('/users/batch-delete', { ids });
};

// 启用/禁用用户
export const toggleUserStatus = (id: number, status: number): Promise<void> => {
  return put<void>(`/users/${id}/status`, { status });
};

// 重置密码
export const resetPassword = (id: number): Promise<{ password: string }> => {
  return post<{ password: string }>(`/users/${id}/reset-password`);
};

// ==================== 角色管理 ====================

// 获取角色列表
export const getRoleList = (params?: PageParams & { keyword?: string; status?: number }): Promise<PageResult<Role>> => {
  return get<PageResult<Role>>('/roles', params as Record<string, unknown>);
};

// 获取所有角色（不分页）
export const getAllRoles = (): Promise<Role[]> => {
  return get<Role[]>('/roles/all');
};

// 获取角色详情
export const getRoleById = (id: number): Promise<Role> => {
  return get<Role>(`/roles/${id}`);
};

// 创建角色
export const createRole = (data: RoleFormData): Promise<Role> => {
  return post<Role>('/roles', data as Record<string, unknown>);
};

// 更新角色
export const updateRole = (id: number, data: RoleFormData): Promise<Role> => {
  return put<Role>(`/roles/${id}`, data as Record<string, unknown>);
};

// 删除角色
export const deleteRole = (id: number): Promise<void> => {
  return del<void>(`/roles/${id}`);
};

// ==================== 权限管理 ====================

// 获取权限树
export const getPermissionTree = (): Promise<Permission[]> => {
  return get<Permission[]>('/permissions/tree');
};

// 获取所有权限
export const getAllPermissions = (): Promise<Permission[]> => {
  return get<Permission[]>('/permissions/all');
};
