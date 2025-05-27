<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { initUserContext } from '$lib/stores/user.svelte';
	import { page } from '$app/state';

	let { children } = $props();
	const { user, login, logout } = initUserContext();
	let hydrated = $state(false); // ✅ new guard

	$effect(() => {
		console.log('🔍 current route:', page.url.pathname);
		console.log('🧠 user.id:', user.id);
		console.log('🧪 hydrated:', hydrated);
	});

	$effect(() => {
		if (!browser || page.url.pathname !== '/') return;

		if (user.id) {
			goto('/dashboard');
		} else {
			goto('/login');
		}
	});

	if (browser) {
		(async () => {
			try {
				const res = await fetch('/api/v1/auth/me', {
					credentials: 'include'
				});
				const result = await res.json();
				console.log('[auth/me result]:', result);

				if (res.ok && result.user?.id) {
					login(result.user);
					console.log('✅ Hydrated user:', result.user);
				} else {
					logout();
					console.warn('❌ No valid user, logging out');
				}
			} catch (err) {
				console.error('❌ Error fetching session:', err);
				logout();
			} finally {
				hydrated = true; // ✅ allow app to continue
			}
		})();
	}
</script>

{#if hydrated}
	{@render children()}
{:else}
	<div class="flex h-screen items-center justify-center text-gray-500">Loading...</div>
{/if}
