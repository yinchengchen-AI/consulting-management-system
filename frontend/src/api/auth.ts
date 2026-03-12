import { post, get } from './request';
import type { 
  LoginParams, 
  LoginResult, 
  UserInfo, 
  ApiResponse 
} from '@/types';

// 登录
export const login = (params: LoginParams): Promise<LoginResult> => {
  return post<LoginResult>('/auth/login', params as Record<string, unknown>);
};

// 登出
export const logout = (): Promise<void> => {
  return post<void>('/auth/logout');
};

// 获取当前用户信息
export const getCurrentUser = (): Promise<UserInfo> => {
  return get<UserInfo>('/auth/current-user');
};

// 刷新token
export const refreshToken = (refreshToken: string): Promise<{ token: string; refreshToken: string; expiresIn: number }> => {
  return post<{ token: string; refreshToken: string; expiresIn: number }>('/auth/refresh-token', { refreshToken });
};

// 修改密码
export const changePassword = (data: { oldPassword: string; newPassword: string }): Promise<void> => {
  return post<void>('/auth/change-password', data as Record<string, unknown>);
};
