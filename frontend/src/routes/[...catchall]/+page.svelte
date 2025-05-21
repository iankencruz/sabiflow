<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { initAuth, isAuthenticated } from '$lib/stores/auth';

	if (browser) {
		(async () => {
			await initAuth();

			if (isAuthenticated) {
				goto('/dashboard');
			} else {
				goto('/login');
			}
		})();
	}
</script>

<div class="flex h-screen flex-col items-center justify-center gap-6 text-gray-700">
	<!-- Animated Spinner -->
	<div
		class="h-12 w-12 animate-spin rounded-full border-4 border-gray-300 border-t-indigo-600"
	></div>

	<!-- Message -->
	<p class="text-base font-medium">Redirecting...</p>
</div>
