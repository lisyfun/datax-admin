// 通用分页响应接口
export interface PaginationResult<T> {
  items: T[];
  total: number;
}

// 用户相关接口
export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  email: string;
  avatar: string;
  status: number;
  roles?: {
    id: number;
    name: string;
    code: string;
  }[];
}

export interface UserListResult extends PaginationResult<UserInfo> {}

// 权限相关接口
export interface Permission {
  id: number;
  parent_id: number;
  name: string;
  code: string;
  type: number;
  status: number;
  children?: Permission[];
}

export interface PermissionTreeResult {
  list: Permission[];
}

// API 响应接口
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}
