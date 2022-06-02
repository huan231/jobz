module.exports = {
  parser: '@typescript-eslint/parser',
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:prettier/recommended',
    'plugin:import/recommended',
    'plugin:import/typescript',
  ],
  plugins: ['@typescript-eslint', 'import', 'prettier', 'react-hooks', 'react'],
  rules: {
    'import/order': [
      'error',
      {
        'groups': ['builtin', 'external', 'internal'],
        'newlines-between': 'always',
      },
    ],
    'import/no-default-export': 'error',
    'react-hooks/exhaustive-deps': 'off',
    'react-hooks/rules-of-hooks': 'error',
  },
  overrides: [
    {
      files: ['*.cjs'],
      env: {
        commonjs: true,
        node: true,
      },
    },
    {
      files: ['vite.config.ts'],
      rules: {
        'import/no-default-export': 'off',
      },
    },
  ],
  ignorePatterns: ['/dist'],
};
