
// src/lib/stores/auth.ts
import { writable, derived } from 'svelte/store';

export type User = {
	id: number;
	firstName: string;
	lastName: string;
	email: string;
	roles?: string[];
};

export const user = writable<User | null>(null);

export const isAuthenticated = derived(user, ($user) => $user !== null);




export async function initAuth() {
	try {
		const res = await fetch('/api/v1/auth/me', { credentials: 'include' });
		console.log('[auth.ts] /api/v1/auth/me status:', res.status);
		const data = await res.json();
		console.log('[auth.ts] /api/v1/auth/me data:', data);
		console.log('[auth.ts] user from server:', data.user);
		user.set(data.user || null);
	} catch (e) {
		console.error('[auth.ts] failed to load user');
		user.set(null);
	}
}


export async function login(credentials: {
	email: string;
	password: string;
}): Promise<{ success: boolean; error?: string }> {
	try {
		const res = await fetch('/api/v1/auth/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(credentials),
			credentials: 'include',
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(error.message || 'Login failed');
		}
		const userData = await res.json();
		user.set(userData);
		return { success: true };
	} catch (err: any) {
		return { success: false, error: err.message };
	}
}

export async function logout() {
	try {
		await fetch('/api/v1/auth/logout', {
			method: 'POST',
			credentials: 'include',
		});
	} catch (err) {
		console.error('Logout failed:', err);
	} finally {
		user.set(null);
		window.location.href = '/login';
	}
}

export function hasRole(role: string): boolean {
	let result = false;

	user.subscribe((val) => {
		result = val?.roles?.includes(role) ?? false;
	})();

	return result;
}

