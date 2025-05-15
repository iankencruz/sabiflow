
// src/lib/stores/userStore.svelte.ts
import type { User } from '$lib/types/user';

let user = $state<User | null>(null);

function setUser(newUser: User) {
	user = newUser;
}

function clearUser() {
	user = null;
}

function getUser() {
	return user;
}

export { getUser, setUser, clearUser };

