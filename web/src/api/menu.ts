import request from '@/utils/request';

export interface CreateMenuRequest {
  parent_id: number;
  name: string;
  path?: string;
  component?: string;
  icon?: string;
  sort?: number;
  hidden?: number;
  cache?: number;
  type: 1 | 2;
  is_external?: number; // 是否为外部链接，0-否，1-是
  external_url?: string; // 外部链接地址
  open_type?: number; // 打开方式，0-内嵌，1-新窗口
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
  hidden: 0 | 1;
  cache: 0 | 1;
  type: 1 | 2;
  is_external?: number; // 是否为外部链接，0-否，1-是
  external_url?: string; // 外部链接地址
  open_type?: number; // 打开方式，0-内嵌，1-新窗口
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
