import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ mode }) => ({
  plugins: [react()],
  build: {
    rollupOptions: {
      input: {
        app: new URL(mode === 'production' ? './index.template.html' : './index.html', import.meta.url).pathname,
      },
    },
  },
}));
