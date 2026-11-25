import { Component } from 'vue';

export interface AppRouteRecordRaw {
  path: string;
  name: string;
  component: Component | (() => Promise<Component>);
  meta: {
    title: string;
    requiresAuth?: boolean;
    icon?: string;
    order?: number;
    hideInMenu?: boolean;
    roles?: string[];
    // 访问该路由所需的权限码（与后端 permissions.code 对应）
    permission?: string;
  };
  children?: AppRouteRecordRaw[];
  redirect?: string;
}
