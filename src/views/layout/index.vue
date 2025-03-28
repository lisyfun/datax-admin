.layout-content {
  padding: 16px;
  overflow: auto;
  background: var(--color-neutral-2);
  transition: all 0.2s ease;
  margin-left: v-bind('collapsed ? "80px" : "220px"');
  margin-top: 64px;
  height: calc(100vh - 64px);
}

/* 当侧边栏收起时的样式 */
.layout-sider.collapsed + .layout-content {
  margin-left: 80px;
}
