<script lang="ts">
	import AuthGuard from '$src/components/AuthGuard.svelte';
</script>

<AuthGuard>
	<h1 class="text-2xl font-bold">Welcome to your dashboard</h1>
	<p class="mt-2 text-gray-600">This content is protected and only visible to logged-in users.</p>
</AuthGuard>
