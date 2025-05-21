<script lang="ts">
	import '$src/app.css';
	import { sidebarNavigation, type NavigationItem } from '$lib/Navigation';
	import { page } from '$app/stores';
	import AuthGuard from '$src/components/AuthGuard.svelte';

	let { children } = $props();
	let isSidebarOpen = $state(false);

	function isActive(path: string): boolean {
		return $page.url.pathname === path || $page.url.pathname.startsWith(path + '/');
	}
</script>

<!-- NavItem Snippet -->
{#snippet NavItem(NavItems: NavigationItem[])}
	{#each NavItems as item}
		<a
			href={item.href}
			class="group flex items-center gap-x-4 rounded-md p-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 hover:text-indigo-600
			{isActive(item.href) ? 'bg-gray-50 text-indigo-600' : ''}"
		>
			<svelte:component this={item.icon} />
			<span>{item.label}</span>
		</a>
	{/each}
{/snippet}

<AuthGuard>
	<div>
		<!-- Mobile sidebar -->
		{#if isSidebarOpen}
			<div class="fixed inset-0 z-50 flex lg:hidden" role="dialog" aria-modal="true">
				<div class="fixed inset-0 bg-gray-900/80" onclick={() => (isSidebarOpen = false)}></div>
				<div class="relative mr-16 flex w-full max-w-xs flex-1">
					<div class="absolute top-0 left-full flex w-16 justify-center pt-5">
						<button class="-m-2.5 p-2.5" onclick={() => (isSidebarOpen = false)}>
							<span class="sr-only">Close sidebar</span>
							<svg
								class="size-6 text-white"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="1.5"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
							</svg>
						</button>
					</div>
					<div class="flex grow flex-col gap-y-5 overflow-y-auto bg-white px-6 pb-2">
						<div class="flex h-16 shrink-0 items-center">
							<img
								class="h-8 w-auto"
								src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
								alt="Logo"
							/>
						</div>
						<nav class="flex flex-1 flex-col">
							<ul class="flex flex-1 flex-col gap-y-7">
								<li>
									<ul class="-mx-2 space-y-1">
										{@render NavItem(sidebarNavigation)}
									</ul>
								</li>
							</ul>
						</nav>
					</div>
				</div>
			</div>
		{/if}
		<!-- Desktop sidebar -->
		<div class="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col">
			<div
				class="flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white px-6"
			>
				<div class="flex h-16 shrink-0 items-center">
					<img
						class="h-8 w-auto"
						src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
						alt="Logo"
					/>
				</div>
				<nav class="flex flex-1 flex-col">
					<ul class="flex flex-1 flex-col gap-y-7">
						<li>
							<ul class="-mx-2 space-y-1">
								{@render NavItem(sidebarNavigation)}
							</ul>
						</li>
					</ul>
				</nav>
			</div>
		</div>
		<!-- Mobile top bar -->
		<div
			class="sticky top-0 z-40 flex items-center gap-x-6 bg-white px-4 py-4 shadow sm:px-6 lg:hidden"
		>
			<button class="-m-2.5 p-2.5 text-gray-700" onclick={() => (isSidebarOpen = true)}>
				<span class="sr-only">Open sidebar</span>
				<svg
					class="size-6"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.5"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
					/>
				</svg>
			</button>
			<div class="flex-1 text-sm font-semibold text-gray-900">Dashboard</div>
			<a href="#">
				<span class="sr-only">Your profile</span>
				<img
					class="size-8 rounded-full bg-gray-50"
					src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
					alt="Profile"
				/>
			</a>
		</div>
		<main class="py-10 lg:pl-72">
			<div class="px-4 sm:px-6 lg:px-8">
				{@render children()}
			</div>
		</main>
	</div>
</AuthGuard>
