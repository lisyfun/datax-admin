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

export interface UpdateMenuRequest extends Partial<CreateMenuRequest> {
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
  status: 0 | 1;
  hidden: boolean;
  cache: boolean;
  type: 1 | 2;
  children?: MenuResponse[];
}

export interface MenuListRequest {
  keyword?: string;
  type?: 1 | 2;
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

// 创建菜单
export function createMenu(data: CreateMenuRequest) {
  return request.post('/menus', data);
}

// 更新菜单
export function updateMenu(id: number, data: UpdateMenuRequest) {
  return request.put(`/menus/${id}`, data);
}

// 删除菜单
export function deleteMenu(id: number) {
  return request.delete(`/menus/${id}`);
}
