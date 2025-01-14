import type { AxiosResponse } from 'axios';
import request from '@/utils/request';
import type {
  Job,
  JobType,
  JobStatus,
  JobShellParams,
  JobHTTPParams,
  JobDataXParams,
  CreateJobRequest,
  UpdateJobRequest,
  JobListRequest,
  JobListResponse,
  JobHistory,
  JobHistoryListRequest,
  JobHistoryListResponse,
} from './types';

export type {
  Job,
  JobType,
  JobStatus,
  JobShellParams,
  JobHTTPParams,
  JobDataXParams,
  CreateJobRequest,
  UpdateJobRequest,
  JobListRequest,
  JobListResponse,
  JobHistory,
  JobHistoryListRequest,
  JobHistoryListResponse,
};

// 获取任务列表
export function getJobList(params: JobListRequest): Promise<AxiosResponse<JobListResponse>> {
  return request.get('/jobs', { params });
}

// 创建任务
export function createJob(data: CreateJobRequest): Promise<AxiosResponse<Job>> {
  console.log(data);
  return request.post('/jobs', data);
}

// 更新任务
export function updateJob(id: number, data: UpdateJobRequest): Promise<AxiosResponse<Job>> {
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

// 获取任务执行历史
export function getJobHistoryList(params: JobHistoryListRequest): Promise<AxiosResponse<JobHistoryListResponse>> {
  return request.get('/jobs/history', { params });
}

// 清理任务历史
export function cleanJobHistory(days: number): Promise<AxiosResponse<void>> {
  return request.post('/jobs/history/clean', { days });
}

