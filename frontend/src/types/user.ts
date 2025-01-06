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
  old_password: string;
  new_password: string;
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
}

export interface UserListParams {
  page: number;
  page_size: number;
  keyword?: string;
}

export interface UserListResult {
  total: number;
  items: UserInfo[];
}

export interface ResetPasswordParams {
  password: string;
}
