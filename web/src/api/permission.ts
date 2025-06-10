import request from '@/utils/request';
import type {
  CreatePermissionParams,
  UpdatePermissionParams,
  PermissionInfo,
  PermissionTreeResult,
} from '@/types/permission';

export function createPermission(data: CreatePermissionParams) {
  return request.post('/permissions', data);
}

export function updatePermission(id: number, data: UpdatePermissionParams) {
  return request.put(`/permissions/${id}`, data);
}

export function deletePermission(id: number) {
  return request.delete(`/permissions/${id}`);
}

export function getPermissionTree() {
  return request.get<PermissionTreeResult>('/permissions');
}

export function getUserPermissions() {
  return request.get('/user/permissions');
}
