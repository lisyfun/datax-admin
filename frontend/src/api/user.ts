import axios from 'axios';
import type {
  LoginParams,
  RegisterParams,
  UpdatePasswordParams,
  UpdateProfileParams,
  UserInfo,
  UserListParams,
  UserListResult,
} from '@/types/user';

export function login(data: LoginParams) {
  return axios.post<{ token: string }>('/api/v1/login', data);
}

export function register(data: RegisterParams) {
  return axios.post('/api/v1/register', data);
}

export function getUserInfo() {
  return axios.get<UserInfo>('/api/v1/user/info');
}

export function updatePassword(data: UpdatePasswordParams) {
  return axios.put('/api/v1/user/password', data);
}

export function updateProfile(data: UpdateProfileParams) {
  return axios.put('/api/v1/user/profile', data);
}

export function getUserList(params: UserListParams) {
  return axios.get<UserListResult>('/api/v1/users', { params });
}

export function updateUserStatus(id: number, status: number) {
  return axios.put(`/api/v1/users/${id}/status`, { status });
}
