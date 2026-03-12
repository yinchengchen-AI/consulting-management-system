import { get, post, put, del } from './request';
import type { 
  Notice,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 通知公告管理 ====================

// 获取通知列表
export const getNoticeList = (params?: PageParams & { 
  keyword?: string; 
  type?: string;
  priority?: string;
  status?: number;
}): Promise<PageResult<Notice>> => {
  return get<PageResult<Notice>>('/notices', params as Record<string, unknown>);
};

// 获取通知详情
export const getNoticeById = (id: number): Promise<Notice> => {
  return get<Notice>(`/notices/${id}`);
};

// 创建通知
export const createNotice = (data: { 
  title: string; 
  content: string; 
  type: 'system' | 'business' | 'announcement';
  priority: 'low' | 'normal' | 'high' | 'urgent';
  publishTime?: string;
  expireTime?: string;
}): Promise<Notice> => {
  return post<Notice>('/notices', data as Record<string, unknown>);
};

// 更新通知
export const updateNotice = (id: number, data: { 
  title: string; 
  content: string; 
  type: 'system' | 'business' | 'announcement';
  priority: 'low' | 'normal' | 'high' | 'urgent';
  publishTime?: string;
  expireTime?: string;
}): Promise<Notice> => {
  return put<Notice>(`/notices/${id}`, data as Record<string, unknown>);
};

// 删除通知
export const deleteNotice = (id: number): Promise<void> => {
  return del<void>(`/notices/${id}`);
};

// 发布通知
export const publishNotice = (id: number): Promise<void> => {
  return put<void>(`/notices/${id}/publish`);
};

// 撤销通知
export const revokeNotice = (id: number): Promise<void> => {
  return put<void>(`/notices/${id}/revoke`);
};

// 标记通知为已读
export const markNoticeAsRead = (id: number): Promise<void> => {
  return put<void>(`/notices/${id}/read`);
};

// 批量标记为已读
export const batchMarkAsRead = (ids: number[]): Promise<void> => {
  return post<void>('/notices/batch-read', { ids });
};

// 获取我的通知列表
export const getMyNotices = (params?: PageParams): Promise<PageResult<Notice>> => {
  return get<PageResult<Notice>>('/notices/my', params as Record<string, unknown>);
};

// 获取未读通知数量
export const getUnreadNoticeCount = (): Promise<number> => {
  return get<number>('/notices/unread-count');
};
