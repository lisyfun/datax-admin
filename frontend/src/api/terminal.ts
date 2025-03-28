import request from '@/utils/request';
import type {
  TerminalInfo,
  CreateTerminalData,
  UpdateTerminalData,
  TerminalListParams,
  TerminalListResponse,
} from '@/types/terminal';
import axios from 'axios';

export default {
  // 创建终端
  createTerminal(data: CreateTerminalData) {
    return request.post<{ message: string }>('/terminals', data);
  },

  // 更新终端
  updateTerminal(id: number, data: UpdateTerminalData) {
    return request.put<{ message: string }>(`/terminals/${id}`, data);
  },

  // 删除终端
  deleteTerminal(id: number) {
    return request.delete<{ message: string }>(`/terminals/${id}`);
  },

  // 获取终端列表
  getTerminalList(params: TerminalListParams) {
    return request.get<TerminalListResponse>('/terminals', { params });
  },

  // 获取终端详情
  getTerminalById(id: number) {
    return request.get<TerminalInfo>(`/terminals/${id}`);
  },

  // 更新终端状态
  updateTerminalStatus(id: number, status: 'online' | 'offline') {
    return request.put<{ message: string }>(`/terminals/${id}/status`, { status });
  },

  // 连接终端
  connectTerminal(id: number) {
    return request.post<{ message: string }>(`/terminals/${id}/connect`);
  },

  // 断开终端连接
  disconnectTerminal(id: number) {
    return request.post<{ message: string }>(`/terminals/${id}/disconnect`);
  },

  // 上传文件到终端
  uploadFiles(id: number, data: FormData) {
    return request.post<{ message: string }>(`/terminals/${id}/upload`, data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      timeout: 30 * 60 * 1000, // 30分钟超时
      maxContentLength: Infinity, // 不限制响应大小
      maxBodyLength: Infinity, // 不限制请求大小
      onUploadProgress: (progressEvent: { loaded: number; total?: number }) => {
        if (!progressEvent.total) return;

        // 获取当前正在上传的文件名
        const formDataFiles = data.getAll('files') as File[];
        if (formDataFiles.length === 0) return;

        // 计算总大小和已上传大小的百分比
        const totalSize = formDataFiles.reduce((acc, file) => acc + file.size, 0);
        const overallProgress = Math.min(100, Math.round((progressEvent.loaded * 100) / totalSize));

        // 找到当前正在上传的文件
        let accumulatedSize = 0;
        let currentFileIndex = 0;

        for (let i = 0; i < formDataFiles.length; i++) {
          const nextSize = accumulatedSize + formDataFiles[i].size;
          if (progressEvent.loaded <= nextSize) {
            currentFileIndex = i;
            break;
          }
          accumulatedSize = nextSize;
        }

        const currentFile = formDataFiles[currentFileIndex];
        const currentFileLoaded = progressEvent.loaded - accumulatedSize;
        const currentFileProgress = Math.min(100, Math.round((currentFileLoaded * 100) / currentFile.size));

        // 触发进度事件
        const event = new CustomEvent('uploadProgress', {
          detail: {
            terminalId: id,
            fileName: currentFile.name,
            progress: currentFileProgress
          }
        });
        window.dispatchEvent(event);
      },
      cancelToken: new axios.CancelToken((cancel) => {
        (window as any).uploadCancel = cancel;
      }),
      validateStatus: (status: number) => {
        return status >= 200 && status < 300;
      }
    });
  },
};
