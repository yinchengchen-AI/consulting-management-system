import request from './request';
import { ApiResponse, PaginatedResponse, Project, PaginationParams } from '@/types';

// 项目查询参数
export interface ProjectListParams extends PaginationParams {
  keyword?: string;
  status?: string;
  type?: string;
  customer_id?: string;
}

// 创建项目参数
export interface CreateProjectParams {
  name: string;
  type?: string;
  description?: string;
  customer_id: string;
  contract_id?: string;
  manager_id?: string;
  start_date?: string;
  end_date?: string;
  budget?: number;
  priority?: number;
}

// 更新项目参数
export interface UpdateProjectParams {
  name?: string;
  type?: string;
  status?: string;
  description?: string;
  manager_id?: string;
  start_date?: string;
  end_date?: string;
  budget?: number;
  actual_cost?: number;
  progress?: number;
  priority?: number;
  deliverables?: string;
  notes?: string;
}

// 获取项目列表
export const getProjectList = (params: ProjectListParams): Promise<ApiResponse<PaginatedResponse<Project>>> => {
  return request.get('/projects', { params });
};

// 获取项目详情
export const getProjectDetail = (id: string): Promise<ApiResponse<Project>> => {
  return request.get(`/projects/${id}`);
};

// 创建项目
export const createProject = (params: CreateProjectParams): Promise<ApiResponse<Project>> => {
  return request.post('/projects', params);
};

// 更新项目
export const updateProject = (id: string, params: UpdateProjectParams): Promise<ApiResponse<Project>> => {
  return request.put(`/projects/${id}`, params);
};

// 删除项目
export const deleteProject = (id: string): Promise<ApiResponse<void>> => {
  return request.delete(`/projects/${id}`);
};
