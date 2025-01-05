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
}

export interface PermissionListParams {
  type?: 'menu' | 'button';
}

export interface PermissionTreeResult {
  list: PermissionInfo[];
}
