<script lang="ts">
	import type { Product } from '$lib/types';
	import { cart } from '$lib/cart.svelte';
	import { priceLabel } from '$lib/format';

	let { product }: { product: Product } = $props();
	const a = $derived(product.attributes);
</script>

<article class="group flex flex-col overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm transition hover:shadow-md">
	<a href="/p/{product.slug}" class="relative block aspect-square bg-gradient-to-b from-slate-50 to-slate-200">
		{#if product.thumb_url}
			<img src={product.thumb_url} alt={product.title} loading="lazy" class="h-full w-full object-contain p-2" />
		{:else}
			<div class="flex h-full items-center justify-center text-4xl text-slate-300">⬢</div>
		{/if}
		{#if product.tier === 'free'}
			<span class="absolute top-2 left-2 rounded-full bg-emerald-500 px-2 py-0.5 text-xs font-semibold text-white">FREE</span>
		{/if}
		{#if product.viewer_url}
			<span class="absolute top-2 right-2 rounded-full bg-white/80 px-2 py-0.5 text-xs text-slate-600">3D</span>
		{/if}
	</a>
	<div class="flex flex-1 flex-col gap-1 p-3">
		<a href="/p/{product.slug}" class="font-medium leading-tight hover:text-indigo-700">{product.title}</a>
		<p class="text-xs text-slate-500">
			{[a.category, a.style, a.polycount ? `${a.polycount.toLocaleString('en')} tris` : null]
				.filter(Boolean)
				.join(' · ')}
		</p>
		<div class="mt-auto flex items-center justify-between pt-2">
			<span class="text-lg font-semibold">{priceLabel(product)}</span>
			<button
				onclick={() => cart.add(product)}
				disabled={cart.has(product.id)}
				class="rounded-lg bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700 disabled:bg-slate-300"
			>
				{cart.has(product.id) ? 'In cart' : 'Add to cart'}
			</button>
		</div>
	</div>
</article>
