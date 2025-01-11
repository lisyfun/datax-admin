import type { AxiosResponse } from 'axios';
import request from '@/utils/request';
import type {
} from './types';

// 获取任务列表
export function getDashboard(): Promise<AxiosResponse<DashboardResponse>> {
  return request.get('/dashboard');
}

