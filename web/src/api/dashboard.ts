import type { AxiosResponse } from 'axios';
import request from '@/utils/request';

interface Stats {
  userCount: number;
  roleCount: number;
  permissionCount: number;
  onlineCount: number;
  successCount: number;
  failedCount: number;
}

interface RecentLogin {
  username: string;
  loginTime: string;
  ip: string;
}

interface JobExecutionTrend {
  date: string;
  successCount: number;
  failedCount: number;
}

interface SystemInfo {
  systemName: string;
  version: string;
  os: string;
  goVersion: string;
  dbVersion: string;
  uptime: string;
  cpuUsage: string;
  memoryTotal: string;
  memoryUsed: string;
  memoryUsage: string;
  diskTotal: string;
  diskUsed: string;
  diskUsage: string;
  numGoroutine: number;
  numCPU: number;
}

interface DashboardResponse {
  stats: Stats;
  recentLogins: RecentLogin[];
  trendData: JobExecutionTrend[];
  systemInfo: SystemInfo;
}

export function getDashboard(): Promise<AxiosResponse<DashboardResponse>> {
  return request.get('/dashboard');
}

