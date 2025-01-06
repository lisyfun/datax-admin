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
  return request.post<{ token: string }>('/v1/login', data);
}

export function register(data: RegisterParams) {
  return request.post('/v1/register', data);
}

export function getUserInfo() {
  return request.get<UserInfo>('/v1/user/info');
}

export function updatePassword(data: UpdatePasswordParams) {
  return request.put('/v1/user/password', data);
}

export function updateProfile(data: UpdateProfileParams) {
  return request.put('/v1/user/profile', data);
}

export function getUserList(params: UserListParams) {
  return request.get<UserListResult>('/v1/users', { params });
}

interface UpdateStatusRequest {
  status: 0 | 1;
}

export function updateUserStatus(id: number, status: number) {
  const data: UpdateStatusRequest = {
    status: status as 0 | 1
  };
  console.log('API call - updateUserStatus:', {
    url: `/v1/users/${id}/status`,
    data,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  });
  return request.put(`/v1/users/${id}/status`, data);
}

export function resetPassword(id: number, data: ResetPasswordParams) {
  return request.put(`/v1/users/${id}/password/reset`, data);
}

export function logout() {
  // 清除本地存储的 token
  localStorage.removeItem('token');
  return Promise.resolve();
}
