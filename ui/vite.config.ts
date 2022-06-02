import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ mode }) => ({
  plugins: [react()],
  build: {
    rollupOptions: {
      input: {
        app: mode === 'production' ? new URL('./index.template.html', import.meta.url).pathname : undefined,
      },
    },
  },
}));
