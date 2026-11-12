<script lang="ts">
	import Chat from '$lib/components/Chat.svelte';
	import Filters from '$lib/components/Filters.svelte';
	import ProductCard from '$lib/components/ProductCard.svelte';
	import type { ActiveFilters, Product } from '$lib/types';

	let { data } = $props();

	let activeFilters = $state<ActiveFilters>({});
	let tier = $state<'' | 'free'>('');
	let text = $state('');
	let products = $state<Product[]>(data.products);
	let first = true;

	// Refetch on every filter change (skip the first run — data already arrived from load).
	$effect(() => {
		const params = new URLSearchParams();
		if (Object.keys(activeFilters).length) params.set('filters', JSON.stringify(activeFilters));
		if (tier) params.set('tier', tier);
		if (text.trim()) params.set('q', text.trim());
		if (first) {
			first = false;
			return;
		}
		fetch(`/api/products?${params}`)
			.then((r) => r.json())
			.then((rows: Product[]) => (products = rows));
	});
</script>

<svelte:head>
	<title>Ol1n Market — 3D models for Roblox Studio</title>
	<meta
		name="description"
		content="Stylized 3D models (FBX + GLB) ready to import into Roblox Studio and other engines. Instant download."
	/>
</svelte:head>

<div class="flex flex-1 flex-col">
	<!-- chat assistant on top, gallery below it (layout from the shop) -->
	<section class="h-[28dvh] min-h-44 border-b-2 border-slate-200 bg-white">
		<Chat />
	</section>

	<section class="flex min-h-0 flex-1 flex-col md:flex-row">
		<aside class="space-y-5 border-b border-slate-200 bg-white p-4 md:w-60 md:shrink-0 md:border-r md:border-b-0">
			<input
				type="search"
				placeholder="Search…"
				bind:value={text}
				class="w-full rounded border border-slate-300 px-2 py-1 text-sm"
			/>
			<label class="flex items-center gap-2 text-sm font-medium text-emerald-700">
				<input type="checkbox" checked={tier === 'free'} onchange={(e) => (tier = e.currentTarget.checked ? 'free' : '')} />
				Free assets only
			</label>
			{#if data.facets.length}
				<Filters schemas={data.facets} bind:active={activeFilters} />
			{/if}
		</aside>

		<div class="flex-1 p-4">
			{#if products.length === 0}
				<p class="mt-8 text-center text-slate-400">
					{data.products.length === 0 ? 'The catalog is being prepared — check back soon.' : 'No assets match these filters.'}
				</p>
			{:else}
				<div class="grid grid-cols-2 gap-4 lg:grid-cols-3 xl:grid-cols-4">
					{#each products as product (product.id)}
						<ProductCard {product} />
					{/each}
				</div>
			{/if}
		</div>
	</section>
</div>
