<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { apiError, priceLabel } from '$lib/format';
	import type { Product } from '$lib/types';

	let { data } = $props();
	let statusFilter = $state<'' | Product['status']>('');
	let busy = $state(false);
	let message = $state('');
	// in-progress edits by product id
	let edits = $state<Record<number, { title: string; price: string; tier: Product['tier'] }>>({});

	const shown = $derived(data.products.filter((p) => !statusFilter || p.status === statusFilter));
	const counts = $derived(
		data.products.reduce<Record<string, number>>((acc, p) => ((acc[p.status] = (acc[p.status] ?? 0) + 1), acc), {})
	);

	function edit(p: Product) {
		return (edits[p.id] ??= { title: p.title, price: (p.price_cents / 100).toFixed(2), tier: p.tier });
	}

	async function save(p: Product, status?: Product['status']) {
		const e = edit(p);
		message = '';
		const res = await fetch(`/api/admin/products/${p.id}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				title: e.title,
				tier: e.tier,
				price_cents: Math.round(parseFloat(e.price || '0') * 100),
				...(status ? { status } : {})
			})
		});
		if (!res.ok) {
			message = `${p.title}: ${await apiError(res)}`;
			return;
		}
		delete edits[p.id];
		await invalidateAll();
	}

	async function runIngest() {
		busy = true;
		message = '';
		const res = await fetch('/api/admin/ingest', { method: 'POST' });
		busy = false;
		if (!res.ok) message = await apiError(res);
		await invalidateAll();
	}
</script>

<svelte:head><title>Admin — Ol1n Market</title></svelte:head>

<div class="mx-auto w-full max-w-6xl space-y-6 p-4 md:p-8">
	<section class="flex flex-wrap items-center gap-3">
		<h1 class="text-2xl font-bold">Admin</h1>
		<span class="text-sm text-slate-500">
			{counts.listed ?? 0} listed · {counts.draft ?? 0} draft · {counts.archived ?? 0} archived ·
			{data.orders.orders.length} objednávek · {data.orders.downloads} stažení
		</span>
		<span class="flex-1"></span>
		{#if data.ingest.enabled}
			<button onclick={runIngest} disabled={busy} class="rounded-lg bg-slate-800 px-3 py-1.5 text-sm text-white disabled:opacity-40">
				{busy ? 'Stahuji z UGC…' : 'Načíst z UGC Factory'}
			</button>
		{:else}
			<span class="text-sm text-amber-700">Ingest vypnutý (UGC_URL)</span>
		{/if}
	</section>

	{#if data.ingest.last}
		{@const l = data.ingest.last}
		<p class="text-xs text-slate-500">
			Poslední ingest {new Date(l.finished_at).toLocaleString('cs')}: nových {l.created}, změněných {l.updated},
			beze změny {l.unchanged}, FAIL přeskočeno {l.skipped_fail}, archivováno {l.archived}
			{#if l.errors.length}<span class="text-red-600"> · chyby: {l.errors.join('; ')}</span>{/if}
		</p>
	{/if}
	{#if message}
		<p class="text-sm text-red-600">{message}</p>
	{/if}

	<section class="space-y-2">
		<div class="flex gap-2 text-sm">
			{#each [['', 'Vše'], ['draft', 'Draft'], ['listed', 'Listed'], ['archived', 'Archived']] as [value, label]}
				<button
					onclick={() => (statusFilter = value as typeof statusFilter)}
					class={'rounded-full px-3 py-1 ' + (statusFilter === value ? 'bg-indigo-600 text-white' : 'bg-white text-slate-600')}
				>
					{label}
				</button>
			{/each}
		</div>

		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white">
			<table class="w-full text-sm">
				<thead class="bg-slate-50 text-left text-xs text-slate-500">
					<tr>
						<th class="p-2"></th>
						<th class="p-2">Název</th>
						<th class="p-2">Info</th>
						<th class="p-2">Tier</th>
						<th class="p-2">Cena (USD)</th>
						<th class="p-2">Stav</th>
						<th class="p-2"></th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100">
					{#each shown as p (p.id)}
						{@const e = edits[p.id] ?? { title: p.title, price: (p.price_cents / 100).toFixed(2), tier: p.tier }}
						<tr>
							<td class="p-2">
								<a href="/p/{p.slug}" target="_blank">
									{#if p.thumb_url}<img src={p.thumb_url} alt="" class="h-12 w-12 rounded bg-slate-100 object-contain" />{/if}
								</a>
							</td>
							<td class="p-2">
								<input value={e.title} oninput={(ev) => (edit(p).title = ev.currentTarget.value)} class="w-56 rounded border border-slate-300 px-2 py-1" />
							</td>
							<td class="p-2 text-xs text-slate-500">
								{p.attributes.category} · {p.attributes.polycount ?? '?'} tris · {p.attributes.verdict}
								{#if p.attributes.ugc_missing}<br /><span class="text-amber-700">zmizel z packed v UGC</span>{/if}
							</td>
							<td class="p-2">
								<select value={e.tier} onchange={(ev) => (edit(p).tier = ev.currentTarget.value as Product['tier'])} class="rounded border border-slate-300 px-1 py-1">
									<option value="paid">paid</option>
									<option value="free">free</option>
								</select>
							</td>
							<td class="p-2">
								<input
									type="number"
									step="0.01"
									min="0"
									value={e.price}
									disabled={e.tier === 'free'}
									oninput={(ev) => (edit(p).price = ev.currentTarget.value)}
									class="w-20 rounded border border-slate-300 px-2 py-1 disabled:opacity-40"
								/>
							</td>
							<td class="p-2">
								<span
									class={'rounded-full px-2 py-0.5 text-xs ' +
										(p.status === 'listed' ? 'bg-emerald-100 text-emerald-800' : p.status === 'draft' ? 'bg-amber-100 text-amber-800' : 'bg-slate-100 text-slate-500')}
								>
									{p.status} · {priceLabel(p)}
								</span>
							</td>
							<td class="space-x-1 p-2 whitespace-nowrap">
								<button onclick={() => save(p)} class="rounded border border-slate-300 px-2 py-1 text-xs hover:bg-slate-50">Uložit</button>
								{#if p.status !== 'listed'}
									<button onclick={() => save(p, 'listed')} class="rounded bg-emerald-600 px-2 py-1 text-xs text-white hover:bg-emerald-700">Vylistovat</button>
								{:else}
									<button onclick={() => save(p, 'draft')} class="rounded bg-slate-200 px-2 py-1 text-xs hover:bg-slate-300">Stáhnout</button>
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</section>

	<section class="space-y-2">
		<h2 class="text-lg font-semibold">Objednávky</h2>
		{#if data.orders.orders.length === 0}
			<p class="text-sm text-slate-500">Zatím žádné.</p>
		{:else}
			<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white">
				<table class="w-full text-sm">
					<thead class="bg-slate-50 text-left text-xs text-slate-500">
						<tr><th class="p-2">#</th><th class="p-2">Kdy</th><th class="p-2">Zákazník</th><th class="p-2">Položek</th><th class="p-2">Celkem</th><th class="p-2">Provider</th><th class="p-2">Stav</th><th class="p-2">Země</th></tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each data.orders.orders as o (o.id)}
							<tr>
								<td class="p-2">{o.id}</td>
								<td class="p-2">{new Date(o.created_at).toLocaleString('cs')}</td>
								<td class="p-2">{o.email}</td>
								<td class="p-2">{o.item_count}</td>
								<td class="p-2">{priceLabel({ tier: o.total_cents ? 'paid' : 'free', price_cents: o.total_cents, currency: o.currency })}</td>
								<td class="p-2">{o.provider}</td>
								<td class="p-2">{o.status}</td>
								<td class="p-2">{o.tax_country ?? '–'}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</section>
</div>
