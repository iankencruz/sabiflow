<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';

	let pingMessage = $state('Connecting...');

	$effect(() => {
		(async () => {
			const res = await fetch(`${PUBLIC_API_URL}/ping`);
			const json = await res.json();
			pingMessage = json.message;
		})();
	});

	console.log(PUBLIC_API_URL); // http://localhost:3000/api
</script>

<h1 class="text-2xl font-bold">{pingMessage}</h1>
