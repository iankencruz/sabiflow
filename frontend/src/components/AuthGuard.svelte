<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { Snippet } from 'svelte';
	import { getUserContext } from '$lib/stores/user.svelte';

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

	const { user, hasRole } = getUserContext();

	let loading = $state(true);
	let authorized = $state(false);

	$effect(() => {
		// If no user logged in
		if (!user.id) {
			authorized = false;

			if (redirect && browser) {
				if (rememberRedirect) {
					sessionStorage.setItem('sabiflow:redirect', page.url.pathname);
				}
				goto(redirectTo);
			}
			return;
		}

		// If role required, check role
		if (requiredRole) {
			if (hasRole(requiredRole)) {
				authorized = true;
			} else {
				authorized = false;
				if (redirect && browser) {
					goto('/unauthorized');
				}
			}
		} else {
			authorized = true;
		}

		loading = false;
	});
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
