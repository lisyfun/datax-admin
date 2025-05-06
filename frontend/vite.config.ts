import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import viteCompression from 'vite-plugin-compression';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    viteCompression({
      verbose: true,
      disable: false,
      threshold: 10240,
      algorithm: 'gzip',
      ext: '.gz',
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/datax/api': {
        target: 'http://localhost:28080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/datax/, ''),
      },
      '/datax/ws': {
        target: 'ws://localhost:28080',
        ws: true,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/datax/, ''),
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    sourcemap: false,
    chunkSizeWarningLimit: 1500,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'arco-vendor': ['@arco-design/web-vue'],
        },
      },
    },
  },
  base: '/datax/',
})
