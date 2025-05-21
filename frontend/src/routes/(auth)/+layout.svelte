<script lang="ts">
	import '$src/app.css';
	import { browser } from '$app/environment';
	import { initAuth } from '$lib/stores/auth';

	let { user, children }: { user: any; children: any } = $props();

	let initializing = $state(true);

	if (browser) {
		(async () => {
			await initAuth();
			initializing = false;
		})();
	}
</script>

{#if user}
	<header class="flex items-center justify-between border-b bg-gray-50 px-6 py-4">
		<h1 class="text-xl font-semibold">Sabiflow</h1>
		<p class="text-sm text-gray-700">Welcome, {user.firstName}</p>
	</header>
{/if}

<main class="p-6">
	{@render children()}
</main>
