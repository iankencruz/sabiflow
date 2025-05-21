<script lang="ts">
	import { browser } from '$app/environment';
	import { initAuth, isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

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
