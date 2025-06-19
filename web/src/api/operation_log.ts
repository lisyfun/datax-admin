import request from '@/utils/request';

// 操作日志接口
export interface OperationLogListParams {
  page: number;
  page_size: number;
  username?: string;
  module?: string;
  action?: string;
  status?: number;
  start_time?: string;
  end_time?: string;
}

export interface OperationLog {
  id: number;
  user_id: number;
  username: string;
  module: string;
  action: string;
  description: string;
  ip: string;
  user_agent: string;
  request_data: string;
  status: number;
  error_msg: string;
  created_at: string;
  updated_at: string;
}

export interface OperationLogListResponse {
  list: OperationLog[];
  total: number;
  page: number;
  page_size: number;
}

export interface BatchDeleteRequest {
  ids: number[];
}

export interface ClearLogsRequest {
  before_days: number;
}

// 获取操作日志列表
export function getOperationLogs(params: OperationLogListParams) {
  return request<OperationLogListResponse>({
    url: '/operation-logs',
    method: 'GET',
    params,
  });
}

// 删除操作日志
export function deleteOperationLog(id: number) {
  return request({
    url: `/operation-logs/${id}`,
    method: 'DELETE',
  });
}

// 批量删除操作日志
export function batchDeleteOperationLogs(data: BatchDeleteRequest) {
  return request({
    url: '/operation-logs/batch-delete',
    method: 'POST',
    data,
  });
}

// 清空操作日志
export function clearOperationLogs(data: ClearLogsRequest) {
  return request({
    url: '/operation-logs/clear',
    method: 'POST',
    data,
  });
}
