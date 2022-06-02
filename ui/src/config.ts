interface Config {
  apiBaseUrl: string;
}

const makeConfig = (): Config => {
  if (import.meta.env.DEV) {
    return { apiBaseUrl: import.meta.env.VITE_API_BASE_URL };
  }

  return { apiBaseUrl: window.APP_CONFIG.apiBaseUrl };
};

export const config = makeConfig();
