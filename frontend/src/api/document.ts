import { get, post, del } from './request';
import type { 
  Document,
  PageResult, 
  PageParams 
} from '@/types';

// ==================== 文档管理 ====================

// 获取文档列表
export const getDocumentList = (params?: PageParams & { 
  keyword?: string; 
  category?: string;
  relatedType?: string;
  relatedId?: number;
}): Promise<PageResult<Document>> => {
  return get<PageResult<Document>>('/documents', params as Record<string, unknown>);
};

// 获取文档详情
export const getDocumentById = (id: number): Promise<Document> => {
  return get<Document>(`/documents/${id}`);
};

// 删除文档
export const deleteDocument = (id: number): Promise<void> => {
  return del<void>(`/documents/${id}`);
};

// 批量删除文档
export const batchDeleteDocuments = (ids: number[]): Promise<void> => {
  return post<void>('/documents/batch-delete', { ids });
};

// 获取文档分类列表
export const getDocumentCategories = (): Promise<string[]> => {
  return get<string[]>('/documents/categories');
};

// 更新文档信息
export const updateDocument = (id: number, data: { name?: string; category?: string; description?: string; tags?: string[] }): Promise<Document> => {
  return post<Document>(`/documents/${id}/update`, data as Record<string, unknown>);
};

// 文件上传（使用FormData）
export const uploadFile = (file: File, data?: { category?: string; description?: string; relatedType?: string; relatedId?: number }): Promise<Document> => {
  const formData = new FormData();
  formData.append('file', file);
  if (data?.category) formData.append('category', data.category);
  if (data?.description) formData.append('description', data.description);
  if (data?.relatedType) formData.append('relatedType', data.relatedType);
  if (data?.relatedId) formData.append('relatedId', data.relatedId.toString());
  
  return post<Document>('/documents/upload', formData as unknown as Record<string, unknown>, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 下载文件
export const downloadFile = (id: number): Promise<Blob> => {
  return get<Blob>(`/documents/${id}/download`, undefined, {
    responseType: 'blob',
  });
};
