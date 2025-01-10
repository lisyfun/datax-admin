import request from '@/utils/request';

export interface CreateMenuRequest {
  parent_id: number;
  name: string;
  path?: string;
  component?: string;
  icon?: string;
  sort?: number;
  hidden?: boolean;
  cache?: boolean;
  type: 1 | 2;
}

export interface UpdateMenuRequest extends CreateMenuRequest {
  status?: 0 | 1;
}

export interface MenuResponse {
  id: number;
  parent_id: number;
  name: string;
  path: string;
  component: string;
  icon: string;
  sort: number;
  status: number;
  hidden: boolean;
  cache: boolean;
  type: number;
  children?: MenuResponse[];
}

export interface MenuListRequest {
  keyword?: string;
  type?: number;
}

export interface MenuListResponse {
  list: MenuResponse[];
}

// 获取用户菜单列表
export function getUserMenus() {
  return request.get<MenuListResponse>('/user/menus');
}

// 获取所有菜单列表（用于菜单管理）
export function getMenuList(params: MenuListRequest = {}) {
  return request.get<MenuListResponse>('/menus', { params });
}

export function createMenu(data: CreateMenuRequest) {
  return request.post<void>('/menus', data);
}

export function updateMenu(id: number, data: UpdateMenuRequest) {
  return request.put<void>(`/menus/${id}`, data);
}

export function deleteMenu(id: number) {
  return request.delete<void>(`/menus/${id}`);
}
