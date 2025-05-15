<script lang="ts">
	import '../app.css';
	import { setContext } from 'svelte';
	import * as userStore from '$lib/stores/userStore.svelte';

	let { children } = $props();
	setContext('userStore', userStore);
</script>

<div class="min-h-screen px-4 py-8">
	{@render children()}
</div>
