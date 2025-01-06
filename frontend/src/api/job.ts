import type { AxiosResponse } from 'axios';
import request from '@/utils/request';

export interface Job {
  id: number;
  name: string;
  description: string;
  cron_expr: string;
  status: number;
  type: number;  // 任务类型：1-Shell脚本 2-DataX任务
  command: string;
  working_dir: string;
  timeout: number;
  retry_times: number;
  retry_interval: number;
  create_time: string;
  update_time: string;
  datax_json?: string;  // DataX 任务的配置 JSON
}

export interface JobListParams {
  page: number;
  page_size: number;
  name?: string;
  status?: number;
}

export interface JobListResponse {
  items: Job[];
  total: number;
}

// 获取任务列表
export function getJobList(params: JobListParams): Promise<AxiosResponse<JobListResponse>> {
  return request.get('/jobs', { params });
}

// 创建任务
export function createJob(data: Partial<Job>): Promise<AxiosResponse<Job>> {
  return request.post('/jobs', data);
}

// 更新任务
export function updateJob(id: number, data: Partial<Job>): Promise<AxiosResponse<Job>> {
  return request.put(`/jobs/${id}`, data);
}

// 删除任务
export function deleteJob(id: number): Promise<AxiosResponse<void>> {
  return request.delete(`/jobs/${id}`);
}

// 启动任务
export function startJob(id: number): Promise<AxiosResponse<void>> {
  return request.post(`/jobs/${id}/start`);
}

// 停止任务
export function stopJob(id: number): Promise<AxiosResponse<void>> {
  return request.post(`/jobs/${id}/stop`);
}

// 执行一次任务
export function executeJob(id: number): Promise<AxiosResponse<void>> {
  return request.post(`/jobs/${id}/execute`);
}

// 批量执行任务
export function executeJobs(ids: number[]): Promise<AxiosResponse<void>> {
  return request.post('/jobs/execute', { ids });
}
