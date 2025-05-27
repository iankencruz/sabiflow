
import { goto } from '$app/navigation';
import { getUserContext } from '$lib/stores/user.svelte';

const BASE_URL = 'http://localhost:8080/api';

export async function fetchApi(endpoint: string, options: RequestInit = {}) {
	const defaultOpts: RequestInit = {
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include'
	};

	const fetchOpts = {
		...defaultOpts,
		...options,
		headers: {
			...defaultOpts.headers,
			...(options.headers || {})
		}
	};

	try {
		const res = await fetch(`${BASE_URL}${endpoint}`, fetchOpts);

		if (res.status === 401) {
			const { logout } = getUserContext();
			logout();
			goto('/login');
			throw new Error('Session expired');
		}

		if (!res.ok) {
			const error = await res.json().catch(() => ({}));
			throw new Error(error.message || `API error: ${res.status}`);
		}

		const contentType = res.headers.get('content-type');
		if (contentType?.includes('application/json')) {
			return res.json();
		}

		return res;
	} catch (err) {
		console.error('API error:', err);
		throw err;
	}
}

