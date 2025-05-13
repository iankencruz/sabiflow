<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { onMount } from 'svelte';
	let name = '';

	onMount(async () => {
		try {
			const res = await fetch(`${PUBLIC_API_URL}/user/me`, {
				credentials: 'include'
			});
			const json = await res.json();
			if (json?.data?.user) {
				const { firstName, lastName } = json.data.user;
				name = `${firstName} ${lastName}`;
			}
		} catch (e) {
			console.error('User not logged in');
		}
	});
</script>

<h1 class="text-2xl font-bold">
	{#if name}
		Hello, {name} 👋
	{:else}
		Welcome to Sabiflow
	{/if}
</h1>
