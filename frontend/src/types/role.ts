export interface CreateRoleParams {
  name: string;
  code: string;
  description?: string;
}

export interface UpdateRoleParams {
  name?: string;
  code?: string;
  description?: string;
  status?: number;
}

export interface RoleInfo {
  id: number;
  name: string;
  code: string;
  description: string;
  status: number;
}

export interface RoleListParams {
  page: number;
  page_size: number;
  keyword?: string;
}

export interface RoleListResult {
  total: number;
  items: RoleInfo[];
}

export interface UpdateRolePermissionsParams {
  permission_ids: number[];
}

export interface UpdateUserRolesParams {
  role_ids: number[];
}
