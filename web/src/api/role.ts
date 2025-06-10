import request from '@/utils/request';
import type {
  CreateRoleParams,
  UpdateRoleParams,
  RoleInfo,
  RoleListParams,
  RoleListResult,
} from '@/types/role';

export function createRole(data: CreateRoleParams) {
  return request.post('/roles', data);
}

export function updateRole(id: number, data: UpdateRoleParams) {
  return request.put(`/roles/${id}`, data);
}

export function deleteRole(id: number) {
  return request.delete(`/roles/${id}`);
}

export function getRoleList(params: RoleListParams) {
  return request.get<RoleListResult>('/roles', { params });
}

export function getRolePermissions(roleId: number) {
  return request.get(`/roles/${roleId}/permissions`);
}

export function updateRolePermissions(roleId: number, permissionIds: number[]) {
  return request.put(`/roles/${roleId}/permissions`, { permission_ids: permissionIds });
}

export function getUserRoles(userId: number) {
  return request.get(`/users/${userId}/roles`);
}

export function updateUserRoles(userId: number, roleIds: number[]) {
  return request.put(`/users/${userId}/roles`, { role_ids: roleIds });
}
