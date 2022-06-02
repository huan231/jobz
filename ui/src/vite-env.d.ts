/// <reference types="vite/client" />

interface ImportMetaEnv {
  VITE_API_BASE_URL: string;
}

interface Window {
  APP_CONFIG: { apiBaseUrl: string };
}
