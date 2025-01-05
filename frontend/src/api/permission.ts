import axios from 'axios';
import type {
  CreatePermissionParams,
  UpdatePermissionParams,
  PermissionInfo,
  PermissionTreeResult,
} from '@/types/permission';

export function createPermission(data: CreatePermissionParams) {
  return axios.post('/api/v1/permissions', data);
}

export function updatePermission(id: number, data: UpdatePermissionParams) {
  return axios.put(`/api/v1/permissions/${id}`, data);
}

export function deletePermission(id: number) {
  return axios.delete(`/api/v1/permissions/${id}`);
}

export function getPermissionTree() {
  return axios.get<PermissionTreeResult>('/api/v1/permissions');
}

export function getUserPermissions() {
  return axios.get('/api/v1/user/permissions');
}
