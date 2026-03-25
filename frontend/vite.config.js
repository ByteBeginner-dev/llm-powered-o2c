/// <reference types="node" />
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';
var __filename = fileURLToPath(import.meta.url);
var __dirname = dirname(__filename);
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            '@': resolve(__dirname, './src'),
        },
    },
    server: {
        port: 3000,
        host: true,
        allowedHosts: 'all',
        proxy: {
            '/api': {
                target: 'https://llm-powered-o2c-production.up.railway.app',
                changeOrigin: true,
                rewrite: function (path) { return path; },
            },
        },
    },
});
