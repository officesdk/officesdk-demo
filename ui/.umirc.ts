import { defineConfig } from "umi";

export default defineConfig({
  routes: [
    { path: "/", redirect: "/showcase" },
    { path: "/showcase", component: "index" },
    { path: "/action", component: "Action/Action" },
    { path: "/file/:id", component: "Action/File", layout: false },
    { path: "/collab/:id", component: "Action/File", layout: false },
  ],
  npmClient: 'yarn',
  mock: {
    include: ['src/mock/**.ts'],
  },
  define: {
    'process.env': process.env
  },
  proxy: {
    '/api': {
      target: 'http://localhost:9301',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
  plugins: ['@umijs/plugins/dist/locale'],
  // 国际化配置
  locale: {
    antd: true,
    default: 'zh-CN',
    baseNavigator: true,
    title: true,
  },
});
