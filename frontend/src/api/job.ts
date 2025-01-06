import axios from 'axios';
import type { TableData } from '@arco-design/web-vue/es/table/interface';

export interface JobRecord {
  id: number;
  name: string;
  type: string;
  description: string;
  cron_expr: string;
  status: number;
  timeout: number;
  retry_count: number;
  retry_delay: number;
  params: any;
  creator: number;
  updater: number;
  created_at: string;
  updated_at: string;
}

export interface JobHistoryRecord {
  id: number;
  job_id: number;
  status: number;
  start_time: string;
  end_time: string;
  duration: number;
  output: string;
  error: string;
  created_at: string;
  updated_at: string;
}

export interface JobParams extends Partial<JobRecord> {
  page: number;
  page_size: number;
}

export interface JobHistoryParams {
  page: number;
  page_size: number;
  job_id?: number;
  status?: number;
  start_time?: string;
  end_time?: string;
}

export function queryJobList(params: JobParams) {
  return axios.get<{
    items: JobRecord[];
    total: number;
  }>('/api/v1/jobs', {
    params,
  });
}

export function queryJobHistoryList(params: JobHistoryParams) {
  return axios.get<{
    items: JobHistoryRecord[];
    total: number;
  }>('/api/v1/jobs/history', {
    params,
  });
}

export function createJob(data: Partial<JobRecord>) {
  return axios.post('/api/v1/jobs', data);
}

export function updateJob(id: number, data: Partial<JobRecord>) {
  return axios.put(`/api/v1/jobs/${id}`, data);
}

export function deleteJob(id: number) {
  return axios.delete(`/api/v1/jobs/${id}`);
}

export function startJob(id: number) {
  return axios.post(`/api/v1/jobs/${id}/start`);
}

export function stopJob(id: number) {
  return axios.post(`/api/v1/jobs/${id}/stop`);
}
