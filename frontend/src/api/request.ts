import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { message } from 'antd';
import { useAuthStore } from '@/stores/auth';

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求队列（用于取消重复请求）
const pendingMap = new Map<string, AbortController>();

// 生成请求key
const getRequestKey = (config: AxiosRequestConfig): string => {
  return `${config.method}_${config.url}_${JSON.stringify(config.params)}_${JSON.stringify(config.data)}`;
};

// 添加请求到队列
const addPending = (config: AxiosRequestConfig): void => {
  const key = getRequestKey(config);
  const controller = new AbortController();
  config.signal = controller.signal;
  pendingMap.set(key, controller);
};

// 移除请求从队列
const removePending = (config: AxiosRequestConfig): void => {
  const key = getRequestKey(config);
  if (pendingMap.has(key)) {
    const controller = pendingMap.get(key);
    controller?.abort();
    pendingMap.delete(key);
  }
};

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 移除重复请求
    removePending(config);
    // 添加请求到队列
    addPending(config);
    
    // 添加token
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 移除请求从队列
    removePending(response.config);
    
    const { data } = response;
    
    // 业务成功
    if (data.code === 200 || data.code === 0) {
      return data.data;
    }
    
    // 业务错误
    message.error(data.message || '请求失败');
    return Promise.reject(new Error(data.message || '请求失败'));
  },
  (error) => {
    // 移除请求从队列
    if (error.config) {
      removePending(error.config);
    }
    
    // 请求被取消
    if (axios.isCancel(error)) {
      return Promise.reject(error);
    }
    
    // 处理HTTP错误
    const { response } = error;
    
    if (response) {
      const { status, data } = response;
      
      switch (status) {
        case 401:
          message.error('登录已过期，请重新登录');
          useAuthStore.getState().logout();
          window.location.href = '/login';
          break;
        case 403:
          message.error('没有权限访问该资源');
          break;
        case 404:
          message.error('请求的资源不存在');
          break;
        case 500:
          message.error('服务器内部错误');
          break;
        default:
          message.error(data?.message || `请求失败: ${status}`);
      }
    } else {
      message.error('网络错误，请检查网络连接');
    }
    
    return Promise.reject(error);
  }
);

// 封装GET请求
export const get = <T>(url: string, params?: Record<string, unknown>, config?: AxiosRequestConfig): Promise<T> => {
  return request.get(url, { params, ...config }) as Promise<T>;
};

// 封装POST请求
export const post = <T>(url: string, data?: Record<string, unknown>, config?: AxiosRequestConfig): Promise<T> => {
  return request.post(url, data, config) as Promise<T>;
};

// 封装PUT请求
export const put = <T>(url: string, data?: Record<string, unknown>, config?: AxiosRequestConfig): Promise<T> => {
  return request.put(url, data, config) as Promise<T>;
};

// 封装DELETE请求
export const del = <T>(url: string, config?: AxiosRequestConfig): Promise<T> => {
  return request.delete(url, config) as Promise<T>;
};

// 封装PATCH请求
export const patch = <T>(url: string, data?: Record<string, unknown>, config?: AxiosRequestConfig): Promise<T> => {
  return request.patch(url, data, config) as Promise<T>;
};

export default request;
