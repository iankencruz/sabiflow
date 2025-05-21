import adapter from '@sveltejs/adapter-static';
import path from 'path';
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
		},
		alias: {
			$components: path.resolve('src/lib/components'),
			$lib: path.resolve('src/lib'),
			$stores: path.resolve('src/lib/stores'),
			$src: path.resolve('src'),
			$assets: path.resolve('src/assets'),
		}
	}
};

export default config;
