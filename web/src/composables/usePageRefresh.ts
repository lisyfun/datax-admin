import { onMounted, onUnmounted } from 'vue';

export function usePageRefresh(refreshCallback: () => void | Promise<void>) {
  // 处理页面刷新
  const handlePageRefresh = () => {
    refreshCallback();
  };

  onMounted(() => {
    // 添加刷新事件监听
    window.addEventListener('page-refresh', handlePageRefresh);
  });

  onUnmounted(() => {
    // 移除刷新事件监听
    window.removeEventListener('page-refresh', handlePageRefresh);
  });

  return {
    refresh: handlePageRefresh
  };
}
