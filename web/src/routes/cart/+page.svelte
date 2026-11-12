<script lang="ts">
	import { cart } from '$lib/cart.svelte';
	import { apiError, priceLabel } from '$lib/format';
	import { goto } from '$app/navigation';

	let { data } = $props();
	let busy = $state(false);
	let message = $state('');

	const currency = $derived(cart.items[0]?.currency ?? 'USD');
	const paidItems = $derived(cart.items.filter((i) => i.tier === 'paid'));

	async function checkout() {
		if (!data.me) {
			goto('/login?next=/cart');
			return;
		}
		busy = true;
		message = '';
		// Phase 1: paid items can't be bought yet — send only the free ones.
		const ids = cart.items.filter((i) => i.tier === 'free').map((i) => i.id);
		const res = await fetch('/api/checkout', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ product_ids: ids })
		});
		busy = false;
		if (!res.ok) {
			message = await apiError(res);
			return;
		}
		const out = await res.json();
		if (out.redirect_url) {
			window.location.href = out.redirect_url;
			return;
		}
		for (const id of ids) cart.remove(id);
		goto('/me?new=1');
	}
</script>

<svelte:head><title>Cart — Ol1n Market</title></svelte:head>

<div class="mx-auto w-full max-w-2xl space-y-4 p-4 md:p-8">
	<h1 class="text-2xl font-bold">Cart</h1>

	{#if cart.items.length === 0}
		<p class="text-slate-500">Your cart is empty. <a href="/" class="underline">Browse assets</a></p>
	{:else}
		<ul class="divide-y divide-slate-200 rounded-xl border border-slate-200 bg-white">
			{#each cart.items as item (item.id)}
				<li class="flex items-center gap-3 p-3">
					{#if item.thumb_url}
						<img src={item.thumb_url} alt="" class="h-14 w-14 rounded bg-slate-100 object-contain" />
					{/if}
					<a href="/p/{item.slug}" class="flex-1 font-medium hover:text-indigo-700">{item.title}</a>
					<span class="font-semibold">{priceLabel(item)}</span>
					<button onclick={() => cart.remove(item.id)} class="text-sm text-slate-400 hover:text-red-600">Remove</button>
				</li>
			{/each}
		</ul>

		<div class="flex items-center justify-between">
			<span class="text-lg">
				Total: <strong>{priceLabel({ tier: cart.totalCents ? 'paid' : 'free', price_cents: cart.totalCents, currency })}</strong>
			</span>
			<button
				onclick={checkout}
				disabled={busy || paidItems.length === cart.items.length}
				class="rounded-lg bg-indigo-600 px-5 py-2 font-medium text-white hover:bg-indigo-700 disabled:opacity-40"
			>
				{data.me ? (paidItems.length ? 'Get the free assets' : 'Get them now') : 'Sign in to continue'}
			</button>
		</div>

		{#if paidItems.length}
			<p class="rounded-lg bg-amber-50 px-3 py-2 text-sm text-amber-800">
				Card and Bitcoin payments are launching soon. Paid assets stay in your cart until then.
			</p>
		{/if}
		{#if message}
			<p class="text-sm text-red-600">{message}</p>
		{/if}
	{/if}
</div>
