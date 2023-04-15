import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 3001,
		host: true,
		proxy: {
			'/api/v1': 'http://localhost:3000',
		}
	}
});
