<script lang="ts">
	import { onMount } from 'svelte';

	// <model-viewer> is vendored from UGC Factory (static/vendor), loaded on the client only.
	let { src, poster, alt }: { src: string; poster?: string; alt: string } = $props();
	let ready = $state(false);

	onMount(() => {
		if (customElements.get('model-viewer')) {
			ready = true;
			return;
		}
		const script = document.createElement('script');
		script.type = 'module';
		script.src = '/vendor/model-viewer.min.js';
		script.onload = () => (ready = true);
		document.head.appendChild(script);
	});
</script>

<div class="relative aspect-square w-full overflow-hidden rounded-2xl bg-gradient-to-b from-slate-100 to-slate-200">
	{#if ready}
		<model-viewer
			{src}
			{poster}
			{alt}
			camera-controls
			auto-rotate
			shadow-intensity="0.6"
			exposure="1.1"
			touch-action="pan-y"
			style="width: 100%; height: 100%;"
		></model-viewer>
	{:else if poster}
		<img src={poster} {alt} class="h-full w-full object-contain" />
	{/if}
	<span class="pointer-events-none absolute bottom-2 left-3 text-xs text-slate-500">Drag to rotate · scroll to zoom</span>
</div>
