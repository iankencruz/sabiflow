import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			// Adjust this to wherever you want the Go backend to serve from
			pages: '../backend/static',
			assets: '../backend/static',
			fallback: 'index.html'
		}),
		paths: {
			base: ''
		}
	}
};

export default config;
