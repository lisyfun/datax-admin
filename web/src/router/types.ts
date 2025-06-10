import { Component } from 'vue';
import type { RouteRecordRaw } from 'vue-router';

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
  };
  children?: AppRouteRecordRaw[];
  redirect?: string;
}
