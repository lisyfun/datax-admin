import axios from 'axios';
import type {
  CreateRoleParams,
  UpdateRoleParams,
  RoleInfo,
  RoleListParams,
  RoleListResult,
} from '@/types/role';

export function createRole(data: CreateRoleParams) {
  return axios.post('/api/v1/roles', data);
}

export function updateRole(id: number, data: UpdateRoleParams) {
  return axios.put(`/api/v1/roles/${id}`, data);
}

export function deleteRole(id: number) {
  return axios.delete(`/api/v1/roles/${id}`);
}

export function getRoleList(params: RoleListParams) {
  return axios.get<RoleListResult>('/api/v1/roles', { params });
}

export function getRolePermissions(roleId: number) {
  return axios.get(`/api/v1/roles/${roleId}/permissions`);
}

export function updateRolePermissions(roleId: number, permissionIds: number[]) {
  return axios.put(`/api/v1/roles/${roleId}/permissions`, { permission_ids: permissionIds });
}

export function getUserRoles(userId: number) {
  return axios.get(`/api/v1/users/${userId}/roles`);
}

export function updateUserRoles(userId: number, roleIds: number[]) {
  return axios.put(`/api/v1/users/${userId}/roles`, { role_ids: roleIds });
}
