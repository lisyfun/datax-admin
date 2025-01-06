import request from '@/utils/request';
import type {
  LoginParams,
  RegisterParams,
  UpdatePasswordParams,
  UpdateProfileParams,
  UserInfo,
  UserListParams,
  UserListResult,
  ResetPasswordParams,
} from '@/types/user';

export function login(data: LoginParams) {
  return request.post<{ token: string }>('/login', data);
}

export function register(data: RegisterParams) {
  return request.post('/register', data);
}

export function getUserInfo() {
  return request.get<UserInfo>('/user/info');
}

export function updatePassword(data: UpdatePasswordParams) {
  return request.put('/user/password', data);
}

export function updateProfile(data: UpdateProfileParams) {
  return request.put('/user/profile', data);
}

export function getUserList(params: UserListParams) {
  return request.get<UserListResult>('/users', { params });
}

interface UpdateStatusRequest {
  status: 0 | 1;
}

export function updateUserStatus(id: number, status: number) {
  const data: UpdateStatusRequest = {
    status: status as 0 | 1
  };
  return request.put(`/users/${id}/status`, data);
}

export function resetPassword(id: number, data: ResetPasswordParams) {
  return request.put(`/users/${id}/password/reset`, data);
}

export function deleteUser(id: number) {
  return request.delete(`/users/${id}`);
}

export function logout() {
  localStorage.removeItem('token');
  return Promise.resolve();
}
