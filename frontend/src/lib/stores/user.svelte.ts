
// src/lib/stores/user.svelte.ts
import { setContext, getContext } from 'svelte';

type User = {
	id: number;
	firstName: string;
	lastName: string;
	email: string;
	roles: string[];
};

const key = Symbol('auth');

export function createAuthStore() {
	// this MUST be declared here and NOT overwritten
	const user = $state<User>({
		id: 0,
		firstName: '',
		lastName: '',
		email: '',
		roles: []
	});

	function login(data: User) {
		user.id = data.id;
		user.firstName = data.firstName;
		user.lastName = data.lastName;
		user.email = data.email;
		user.roles = data.roles;
	}

	function logout() {
		user.id = 0;
		user.firstName = '';
		user.lastName = '';
		user.email = '';
		user.roles = [];
	}

	function hasRole(role: string) {
		return user.roles.includes(role);
	}

	return {
		user: user,
		login,
		logout,
		hasRole
	};
}

export function initUserContext() {
	const store = createAuthStore();
	setContext(key, store);
	return store;
}

export function getUserContext() {
	return getContext<ReturnType<typeof createAuthStore>>(key);
}

