export type JobType = 'shell' | 'http' | 'datax';
export type JobStatus = 0 | 1 | 2; // 0:停止 1:运行中 2:禁用

export interface JobShellParams {
  command: string;
  work_dir: string;
  environment: Record<string, string>;
}

export interface JobHTTPParams {
  url: string;
  method: string;
  headers: Record<string, string>;
  body: string;
  success_code: number[];
}

export interface JobDataXParams {
  job_path: string;
  parameters: Record<string, string>;
  jvm_options: string[];
  speed: number;
}

export interface Job {
  id: number;
  name: string;
  type: JobType;
  description: string;
  cron_expr: string;
  status: JobStatus;
  timeout: number;
  retry_count: number;
  retry_delay: number;
  params: string;
  creator: number;
  updater?: number;
  created_at: string;
  updated_at: string;
}

export interface CreateJobRequest {
  name: string;
  type: JobType;
  description?: string;
  cron_expr: string;
  timeout?: number;
  retry_count?: number;
  retry_delay?: number;
  params: JobShellParams | JobHTTPParams | JobDataXParams;
}

export interface UpdateJobRequest {
  name?: string;
  description?: string;
  cron_expr?: string;
  timeout?: number;
  retry_count?: number;
  retry_delay?: number;
  params?: JobShellParams | JobHTTPParams | JobDataXParams;
}

export interface JobListRequest {
  page: number;
  page_size: number;
  keyword?: string;
  type?: JobType;
  status?: JobStatus;
}

export interface JobListResponse {
  items: Job[];
  total: number;
}
