<script lang="ts">
	import { page } from '$app/state';
	import { apiError } from '$lib/format';

	let { data } = $props();
	let error = $state('');

	// The download token is minted on every click (24 h, 5 uses).
	async function download(id: number) {
		error = '';
		const res = await fetch(`/api/me/licenses/${id}/download`, { method: 'POST' });
		if (!res.ok) {
			error = await apiError(res);
			return;
		}
		const { url } = await res.json();
		window.location.href = url;
	}
</script>

<svelte:head><title>My assets — Ol1n Market</title></svelte:head>

<div class="mx-auto w-full max-w-3xl space-y-4 p-4 md:p-8">
	<h1 class="text-2xl font-bold">My assets</h1>
	<p class="text-sm text-slate-500">Signed in as {data.me?.email}</p>

	{#if page.url.searchParams.get('new')}
		<p class="rounded-lg bg-emerald-50 px-3 py-2 text-emerald-800">Done! Your assets are ready to download below.</p>
	{/if}

	{#if data.licenses.length === 0}
		<p class="text-slate-500">Nothing here yet. <a href="/" class="underline">Browse assets</a></p>
	{:else}
		<ul class="divide-y divide-slate-200 rounded-xl border border-slate-200 bg-white">
			{#each data.licenses as lic (lic.id)}
				<li class="flex flex-wrap items-center gap-3 p-3">
					<div class="flex-1">
						<a href="/p/{lic.slug}" class="font-medium hover:text-indigo-700">{lic.title}</a>
						<p class="text-xs text-slate-500">
							{lic.license_type} license · <span class="font-mono">{lic.license_key}</span> · {new Date(lic.issued_at).toLocaleDateString('en')}
						</p>
					</div>
					<button onclick={() => download(lic.id)} class="rounded-lg bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700">
						Download .zip
					</button>
				</li>
			{/each}
		</ul>
	{/if}
	{#if error}
		<p class="text-sm text-red-600">{error}</p>
	{/if}
</div>
