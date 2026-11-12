<script lang="ts">
	import ModelViewer from '$lib/components/ModelViewer.svelte';
	import { cart } from '$lib/cart.svelte';
	import { bytesLabel, priceLabel } from '$lib/format';
	import { goto } from '$app/navigation';

	let { data } = $props();
	const p = $derived(data.product);
	const a = $derived(p.attributes);

	const FORMAT_INFO: Record<string, string> = {
		fbx: 'FBX — Roblox-ready mesh (≤ 4,000 tris, scaled in studs)',
		png: 'PNG — baked texture for the FBX',
		glb: 'GLB — higher-detail source mesh for other engines',
		obj: 'OBJ — universal mesh',
		rbxm: 'RBXM — Roblox model file'
	};

	const specs = $derived(
		[
			['Category', a.category],
			['Style', a.style],
			['Collection', a.collection],
			['Triangles', a.polycount?.toLocaleString('en')],
			['Texture', a.texture_px ? `${a.texture_px} px` : undefined],
			['Attachment', a.roblox_attachment],
			['Size (studs)', a.bbox_studs?.map((n) => n.toFixed(2)).join(' × ')],
			['Watertight', a.watertight === undefined ? undefined : a.watertight ? 'yes' : 'no']
		].filter(([, v]) => v) as [string, string][]
	);

	function getIt() {
		cart.add(p);
		goto('/cart');
	}
</script>

<svelte:head>
	<title>{p.title} — Ol1n Market</title>
	<meta name="description" content={p.description || `${p.title} — 3D model for Roblox Studio (FBX + GLB).`} />
</svelte:head>

<div class="mx-auto grid w-full max-w-6xl gap-8 p-4 md:grid-cols-[3fr_2fr] md:p-8">
	<div>
		{#if p.viewer_url}
			<ModelViewer src={p.viewer_url} poster={p.thumb_url} alt={p.title} />
		{:else if p.thumb_url}
			<img src={p.thumb_url} alt={p.title} class="aspect-square w-full rounded-2xl bg-slate-100 object-contain" />
		{/if}
		{#if p.thumb_url && p.viewer_url}
			<div class="mt-3 flex items-center gap-3 text-xs text-slate-500">
				<img src={p.thumb_url} alt="Concept art" class="h-16 w-16 rounded-lg bg-slate-100 object-contain" />
				Concept art the model was generated from.
			</div>
		{/if}
	</div>

	<div class="space-y-5">
		{#if p.status !== 'listed'}
			<p class="rounded-lg bg-amber-100 px-3 py-2 text-sm text-amber-800">Admin preview — status: {p.status}</p>
		{/if}
		<div>
			<h1 class="text-2xl font-bold tracking-tight">{p.title}</h1>
			{#if a.tags?.length}
				<div class="mt-2 flex flex-wrap gap-1.5">
					{#each a.tags as tag}
						<span class="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-600">{tag}</span>
					{/each}
				</div>
			{/if}
		</div>

		<div class="flex items-center gap-4">
			<span class="text-3xl font-semibold">{priceLabel(p)}</span>
			{#if p.owned}
				<a href="/me" class="rounded-lg bg-emerald-600 px-4 py-2 font-medium text-white hover:bg-emerald-700">
					You own this — download
				</a>
			{:else if cart.has(p.id)}
				<a href="/cart" class="rounded-lg bg-slate-800 px-4 py-2 font-medium text-white">In cart — checkout</a>
			{:else}
				<button onclick={getIt} class="rounded-lg bg-indigo-600 px-4 py-2 font-medium text-white hover:bg-indigo-700">
					{p.tier === 'free' ? 'Get it free' : 'Add to cart'}
				</button>
			{/if}
		</div>
		<p class="text-xs text-slate-500">
			Instant download · {p.tier === 'free' ? 'Personal license' : 'Personal license included'} ·
			<a href="/licenses" class="underline">license terms</a>
		</p>

		{#if p.description}
			<p class="whitespace-pre-line text-slate-700">{p.description}</p>
		{/if}

		<div>
			<h2 class="mb-2 text-sm font-semibold text-slate-700">What's in the download</h2>
			<ul class="space-y-1 text-sm text-slate-600">
				{#each p.files as f}
					<li>{FORMAT_INFO[f.format] ?? f.format.toUpperCase()} <span class="text-slate-400">({bytesLabel(f.bytes)})</span></li>
				{/each}
				<li>LICENSE.txt + README.txt (import guide for Roblox Studio)</li>
			</ul>
		</div>

		{#if specs.length}
			<dl class="grid grid-cols-2 gap-x-4 gap-y-1 rounded-xl border border-slate-200 bg-white p-4 text-sm">
				{#each specs as [k, v]}
					<dt class="text-slate-500">{k}</dt>
					<dd class="font-medium">{v}</dd>
				{/each}
			</dl>
		{/if}

		<p class="text-xs text-slate-400">
			Made with AI-assisted 3D generation and cleanup. Uploads to Roblox go through Roblox moderation, which we don't control —
			see the refund policy.
		</p>
	</div>
</div>
