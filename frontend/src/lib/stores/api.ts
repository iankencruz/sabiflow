
// src/lib/api.ts
import { user } from '$lib/stores/auth';
import { goto } from '$app/navigation';

const BASE_URL = 'http://localhost:8080/api';

export async function fetchApi(endpoint: string, options: RequestInit = {}) {
	const defaultOpts: RequestInit = {
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
	};

	const fetchOpts: RequestInit = {
		...defaultOpts,
		...options,
		headers: {
			...defaultOpts.headers,
			...options.headers,
		},
	};

	try {
		const res = await fetch(`${BASE_URL}${endpoint}`, fetchOpts);

		if (res.status === 401) {
			user.set(null);
			goto('/login?reason=session_expired');
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
		console.error('Fetch failed:', err);
		throw err;
	}
}
