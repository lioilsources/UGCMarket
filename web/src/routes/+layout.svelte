<script lang="ts">
	import '../app.css';
	import { cart } from '$lib/cart.svelte';
	import { invalidateAll } from '$app/navigation';

	let { children, data } = $props();

	async function logout() {
		await fetch('/api/auth/logout', { method: 'POST' });
		await invalidateAll();
	}
</script>

<div class="flex min-h-dvh flex-col">
	<header class="sticky top-0 z-10 flex items-center justify-between border-b border-slate-200 bg-white/95 px-4 py-2 backdrop-blur">
		<a href="/" class="flex items-center gap-2 text-lg font-bold tracking-tight">
			<img src="/favicon.svg" alt="" class="h-7 w-7" />
			Ol1n Market
		</a>
		<nav class="flex items-center gap-3 text-sm">
			{#if data.me?.is_admin}
				<a href="/admin" class="text-slate-500 hover:text-slate-900">Admin</a>
			{/if}
			{#if data.me}
				<a href="/me" class="text-slate-600 hover:text-slate-900">My assets</a>
				<button onclick={logout} class="text-slate-400 hover:text-slate-700">Sign out</button>
			{:else}
				<a href="/login" class="text-slate-600 hover:text-slate-900">Sign in</a>
			{/if}
			<a href="/cart" class="rounded-full bg-indigo-950 px-3 py-1 font-medium text-white">
				Cart · {cart.count}
			</a>
		</nav>
	</header>
	<main class="flex min-h-0 flex-1 flex-col">
		{@render children()}
	</main>
	<footer class="border-t border-slate-200 bg-white px-4 py-3 text-center text-xs text-slate-500">
		Ol1n Market sells independent 3D models compatible with Roblox Studio. Not affiliated with, endorsed by, or
		sponsored by Roblox Corporation. “Roblox” is a trademark of Roblox Corporation.
		<a href="/licenses" class="underline">Licenses &amp; refunds</a>
	</footer>
</div>
