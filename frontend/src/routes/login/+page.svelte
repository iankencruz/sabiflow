<script lang="ts">
	import { goto } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	let email = $state('');
	let password = $state('');
	let error = $state('');

	async function handleLogin(e: SubmitEvent) {
		e.preventDefault();

		const res = await fetch('/api/auth/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password })
		});

		const data = await res.json();

		if (!res.ok) {
			error = data.message || 'Login failed';
			return;
		}

		goto('/');
	}
</script>

<form onsubmit={handleLogin} class="mx-auto mt-20 max-w-md space-y-4">
	<h1 class="text-center text-2xl font-bold">Login</h1>

	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}

	<div>
		<label for="email">Email</label>
		<input id="email" type="email" bind:value={email} class="w-full rounded border p-2" required />
	</div>

	<div>
		<label for="password">Password</label>
		<input
			id="password"
			type="password"
			bind:value={password}
			class="w-full rounded border p-2"
			required
		/>
	</div>

	<button type="submit" class="w-full rounded bg-black p-2 text-white">Login</button>
</form>
