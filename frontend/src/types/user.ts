import type { RoleInfo } from './role';

export interface LoginParams {
  username: string;
  password: string;
}

export interface RegisterParams {
  username: string;
  password: string;
  nickname?: string;
  email?: string;
}

export interface UpdatePasswordParams {
  oldPassword: string;
  newPassword: string;
}

export interface UpdateProfileParams {
  nickname?: string;
  email?: string;
  avatar?: string;
}

export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  email: string;
  avatar: string;
  status: number;
  roles?: RoleInfo[];
}

export interface UserListParams {
  page: number;
  page_size: number;
  keyword?: string;
}

export interface UserListResult {
  total: number;
  list: UserInfo[];
}

export interface ResetPasswordParams {
  password: string;
}

export interface UserForm {
  id?: number;
  username: string;
  password?: string;
  nickname?: string;
  email?: string;
}

export interface PasswordForm {
  oldPassword: string;
  newPassword: string;
}

export interface ResetPasswordForm {
  password: string;
}
