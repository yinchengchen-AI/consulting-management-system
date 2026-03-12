import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { UserInfo } from '@/types';

interface AuthState {
  // 状态
  token: string | null;
  refreshToken: string | null;
  userInfo: UserInfo | null;
  isAuthenticated: boolean;
  
  // 方法
  setToken: (token: string, refreshToken: string) => void;
  setUserInfo: (userInfo: UserInfo) => void;
  logout: () => void;
  hasPermission: (permission: string) => boolean;
  hasRole: (role: string) => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // 初始状态
      token: null,
      refreshToken: null,
      userInfo: null,
      isAuthenticated: false,
      
      // 设置token
      setToken: (token: string, refreshToken: string) => {
        set({ 
          token, 
          refreshToken, 
          isAuthenticated: true 
        });
      },
      
      // 设置用户信息
      setUserInfo: (userInfo: UserInfo) => {
        set({ userInfo });
      },
      
      // 登出
      logout: () => {
        set({ 
          token: null, 
          refreshToken: null, 
          userInfo: null, 
          isAuthenticated: false 
        });
      },
      
      // 检查是否有权限
      hasPermission: (permission: string): boolean => {
        const { userInfo } = get();
        if (!userInfo?.permissions) return false;
        return userInfo.permissions.includes(permission) || userInfo.permissions.includes('*');
      },
      
      // 检查是否有角色
      hasRole: (role: string): boolean => {
        const { userInfo } = get();
        if (!userInfo?.roles) return false;
        return userInfo.roles.includes(role) || userInfo.roles.includes('admin');
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ 
        token: state.token, 
        refreshToken: state.refreshToken,
        userInfo: state.userInfo,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
