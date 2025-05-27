<script lang="ts">
	import '$src/app.css';
	import { browser } from '$app/environment';
	import type { LoadEvent } from '@sveltejs/kit';

	import { getUserContext } from '$lib/stores/user.svelte';

	let { children }: { user: any; children: any } = $props();

	const { user } = getUserContext();
</script>

<main>
	{@render children()}
</main>
