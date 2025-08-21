import { TreeNodeData } from '@arco-design/web-vue';
import { VNode, h } from 'vue';

export interface PermissionInfo {
  id: number;
  name: string;
  code: string;
  type: 'menu' | 'button';
  parent_id: number;
  path: string;
  component: string;
  icon: string;
  sort: number;
  status: number;
  hidden: number;
  cache: number;
  is_external?: number; // 是否为外部链接，0-否，1-是
  external_url?: string; // 外部链接地址
  open_type?: number; // 打开方式，0-内嵌，1-新窗口
  children?: PermissionInfo[];
}

export const convertToTreeData = (permissions: PermissionInfo[]): TreeNodeData[] => {
  return permissions.map(item => ({
    key: item.id,
    title: item.name,
    selectable: true,
    disabled: item.status === 0,
    icon: item.icon ? () => h('i', { class: item.icon }) as VNode : undefined,
    children: item.children ? convertToTreeData(item.children) : undefined,
    isLeaf: !item.children || item.children.length === 0
  }));
};

export interface CreatePermissionParams {
  name: string;
  code: string;
  type: 'menu' | 'button';
  parent_id?: number;
  path?: string;
  component?: string;
  icon?: string;
  sort: number;
  hidden: number;
  cache: number;
  is_external?: number; // 是否为外部链接，0-否，1-是
  external_url?: string; // 外部链接地址
  open_type?: number; // 打开方式，0-内嵌，1-新窗口
}

export interface UpdatePermissionParams {
  name?: string;
  code?: string;
  type?: 'menu' | 'button';
  parent_id?: number;
  path?: string;
  component?: string;
  icon?: string;
  sort?: number;
  status?: number;
  hidden?: number;
  cache?: number;
  is_external?: number; // 是否为外部链接，0-否，1-是
  external_url?: string; // 外部链接地址
  open_type?: number; // 打开方式，0-内嵌，1-新窗口
}

export interface PermissionListParams {
  type?: 'menu' | 'button';
}

export interface PermissionTreeResult {
  list: PermissionInfo[];
}
