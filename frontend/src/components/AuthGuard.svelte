<!-- src/lib/components/AuthGuard.svelte -->

<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { isAuthenticated, initAuth, hasRole } from '$lib/stores/auth';
	import type { Snippet } from 'svelte';
	import { browser } from '$app/environment';

	interface Props {
		requiredRole?: string;
		redirect?: boolean;
		redirectTo?: string;
		rememberRedirect?: boolean;
		children: Snippet;
	}

	let {
		requiredRole = undefined,
		redirect = true,
		redirectTo = '/login',
		rememberRedirect = false,
		children
	}: Props = $props();

	let loading = $state(true);
	let authorized = $state(false);

	// ✅ Fetch session info from /api/auth/me

	(async () => {
		await initAuth();
		checkAuth();
		loading = false;
	})();

	function checkAuth() {
		if (!$isAuthenticated) {
			authorized = false;

			if (redirect && browser) {
				goto('/login');
			}
			return;
		}

		if (requiredRole) {
			authorized = hasRole(requiredRole);
			if (!authorized && redirect) {
				goto('/unauthorized');
			}
		} else {
			authorized = true;
		}
	}
</script>

{#if loading}
	<div class="flex h-screen items-center justify-center text-lg font-medium">Loading...</div>
{:else if authorized}
	{@render children()}
{:else if !redirect}
	<div class="mt-10 text-center">
		<h2 class="text-xl font-semibold">Not Authorized</h2>
		<p class="text-gray-600">You don’t have permission to view this content.</p>
	</div>
{/if}
