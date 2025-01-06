import request from '@/utils/request';
import type {
  CreateRoleParams,
  UpdateRoleParams,
  RoleInfo,
  RoleListParams,
  RoleListResult,
} from '@/types/role';

export function createRole(data: CreateRoleParams) {
  return request.post('/v1/roles', data);
}

export function updateRole(id: number, data: UpdateRoleParams) {
  return request.put(`/v1/roles/${id}`, data);
}

export function deleteRole(id: number) {
  return request.delete(`/v1/roles/${id}`);
}

export function getRoleList(params: RoleListParams) {
  return request.get<RoleListResult>('/v1/roles', { params });
}

export function getRolePermissions(roleId: number) {
  return request.get(`/v1/roles/${roleId}/permissions`);
}

export function updateRolePermissions(roleId: number, permissionIds: number[]) {
  return request.put(`/v1/roles/${roleId}/permissions`, { permission_ids: permissionIds });
}

export function getUserRoles(userId: number) {
  return request.get(`/v1/users/${userId}/roles`);
}

export function updateUserRoles(userId: number, roleIds: number[]) {
  return request.put(`/v1/users/${userId}/roles`, { role_ids: roleIds });
}
