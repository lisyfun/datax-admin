import request from '@/utils/request';
import type {
  CreatePermissionParams,
  UpdatePermissionParams,
  PermissionInfo,
  PermissionTreeResult,
} from '@/types/permission';

export function createPermission(data: CreatePermissionParams) {
  return request.post('/v1/permissions', data);
}

export function updatePermission(id: number, data: UpdatePermissionParams) {
  return request.put(`/v1/permissions/${id}`, data);
}

export function deletePermission(id: number) {
  return request.delete(`/v1/permissions/${id}`);
}

export function getPermissionTree() {
  return request.get<PermissionTreeResult>('/v1/permissions');
}

export function getUserPermissions() {
  return request.get('/v1/user/permissions');
}
