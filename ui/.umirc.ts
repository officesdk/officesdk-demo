import { defineConfig } from "umi";

export default defineConfig({
  routes: [
    { path: "/", redirect: "/showcase" },
    { path: "/showcase", component: "index" },
    { path: "/action", component: "Action/Action" },
    { path: "/case", component: "Case/Case", layout: false },
    { path: "/file/:id", component: "Action/File", layout: false },
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
      target: 'http://172.16.21.165:9301',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
});
