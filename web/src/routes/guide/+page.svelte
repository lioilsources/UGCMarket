<script lang="ts">
	import { GUIDE } from '$lib/guide/guide.js';
	import { renderGuide } from '$lib/guide/render.js';

	// Unverified sections show up only in `vite dev`, marked as draft.
	const g = renderGuide(GUIDE, { drafts: import.meta.env.DEV });

	let zoomed = $state<string | null>(null);
	function onclick(e: MouseEvent) {
		const t = e.target as HTMLElement;
		if (t instanceof HTMLImageElement && t.hasAttribute('data-zoom')) zoomed = t.src;
	}
</script>

<svelte:head>
	<title>Studio guide — Ol1n Market</title>
	<meta name="description" content="How to use a model from Ol1n Market in Roblox Studio: import, texture, accessory, try it on." />
</svelte:head>

<svelte:window onkeydown={(e) => e.key === 'Escape' && (zoomed = null)} />

<article class="guide mx-auto w-full max-w-3xl space-y-6 p-4 md:p-8">
	<header class="space-y-2">
		<h1 class="text-2xl font-bold">From ZIP to your character in Roblox Studio</h1>
		<p class="text-slate-600">
			Step by step: what is in the download, how to bring the model into Roblox Studio and how to wear it.
		</p>
		{#if GUIDE.studioVersion}
			<p class="text-xs text-slate-500">Checked in Roblox Studio {GUIDE.studioVersion} on {GUIDE.verifiedOn}.</p>
		{/if}
	</header>

	{#if !g.complete && !import.meta.env.DEV}
		<p class="rounded-lg bg-amber-50 px-3 py-2 text-sm text-amber-800">
			The step-by-step Studio walkthrough with screenshots is being written. This page grows as each step is tested.
		</p>
	{/if}

	<nav>{@html g.toc}</nav>

	<!-- svelte-ignore a11y_click_events_have_key_events, a11y_no_static_element_interactions -->
	<div {onclick}>{@html g.body}</div>
</article>

{#if zoomed}
	<button class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4" onclick={() => (zoomed = null)} aria-label="Close image">
		<img src={zoomed} alt="" class="max-h-full max-w-full rounded" />
	</button>
{/if}

<style>
	.guide :global(h2) { margin: 2rem 0 0.75rem; font-size: 1.25rem; font-weight: 700; scroll-margin-top: 1rem; }
	.guide :global(p) { margin: 0.5rem 0; line-height: 1.6; }
	.guide :global(code) { border-radius: 0.25rem; background: #f1f5f9; padding: 0 0.25rem; font-size: 0.9em; }
	.guide :global(pre) { overflow-x: auto; border-radius: 0.5rem; background: #0f172a; padding: 0.75rem; color: #e2e8f0; }
	.guide :global(pre code) { background: none; padding: 0; }
	.guide :global(ol.steps) { list-style: decimal; padding-left: 1.5rem; }
	.guide :global(ol.steps li) { margin: 0.75rem 0; line-height: 1.6; }
	.guide :global(ol.toc) { list-style: decimal; padding-left: 1.5rem; color: #4f46e5; }
	.guide :global(img[data-zoom]) { display: block; margin-top: 0.5rem; max-width: 100%; height: auto; cursor: zoom-in; border: 1px solid #e2e8f0; border-radius: 0.5rem; }
	.guide :global(figcaption) { font-size: 0.8rem; color: #64748b; }
	.guide :global(.table) { overflow-x: auto; }
	.guide :global(table) { width: 100%; border-collapse: collapse; font-size: 0.9rem; }
	.guide :global(th), .guide :global(td) { border-bottom: 1px solid #e2e8f0; padding: 0.4rem 0.5rem; text-align: left; vertical-align: top; }
	.guide :global(tr:target) { background: #eef2ff; }
	.guide :global(.draft) { border-radius: 0.25rem; background: #fef3c7; padding: 0 0.4rem; font-size: 0.7rem; color: #92400e; vertical-align: middle; }
	.guide :global(ul.todo) { border-left: 3px solid #f59e0b; padding-left: 0.75rem; font-size: 0.85rem; color: #92400e; }
</style>
