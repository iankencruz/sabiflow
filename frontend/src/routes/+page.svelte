<script lang="ts">
	import { getContext } from 'svelte';
	import type { User } from '$lib/types/user';

	let { user, setUser, clearUser } = getContext('userStore') as {
		user: User | null;
		setUser: (user: User) => void;
		clearUser: () => void;
	};

	let name: string | null = null;
	if (user) {
		name = `${user.firstName} ${user.lastName}`.trim();
	}
</script>

<h1 class="mt-10 text-center text-2xl font-bold">
	{#if name}
		Welcome, {name}
	{:else}
		Welcome to Sabiflow
	{/if}
</h1>
