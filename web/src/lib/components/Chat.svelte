<script lang="ts">
	import type { ChatMsg, Product } from '$lib/types';
	import { priceLabel } from '$lib/format';

	let messages = $state<ChatMsg[]>([]);
	let input = $state('');
	let busy = $state(false);
	let scroller: HTMLDivElement;

	function scrollDown() {
		requestAnimationFrame(() => scroller?.scrollTo({ top: scroller.scrollHeight }));
	}

	async function send(e: SubmitEvent) {
		e.preventDefault();
		const question = input.trim();
		if (!question || busy) return;
		input = '';
		busy = true;
		messages.push({ role: 'user', content: question });
		messages.push({ role: 'assistant', content: '', products: [] });
		// the proxy from the array must be mutated, not the original object — otherwise changes don't render
		const assistant = messages[messages.length - 1];
		scrollDown();

		try {
			const res = await fetch('/api/chat', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					messages: messages
						.slice(0, -1)
						.map(({ role, content }) => ({ role, content }))
				})
			});
			if (!res.ok || !res.body) throw new Error(`chat failed (${res.status})`);

			const reader = res.body.pipeThrough(new TextDecoderStream()).getReader();
			let buffer = '';
			for (;;) {
				const { done, value } = await reader.read();
				if (done) break;
				buffer += value;
				let sep;
				while ((sep = buffer.indexOf('\n\n')) !== -1) {
					const block = buffer.slice(0, sep);
					buffer = buffer.slice(sep + 2);
					const event = /^event: (.+)$/m.exec(block)?.[1];
					const dataRaw = /^data: (.+)$/m.exec(block)?.[1];
					if (!event || !dataRaw) continue;
					const data = JSON.parse(dataRaw);
					if (event === 'token') assistant.content += data.text;
					else if (event === 'products') assistant.products = data as Product[];
					scrollDown();
				}
			}
		} catch (err) {
			assistant.content ||= `Couldn't reach the assistant: ${err}`;
		} finally {
			busy = false;
			scrollDown();
		}
	}
</script>

<div class="flex h-full flex-col">
	<div bind:this={scroller} class="flex-1 space-y-3 overflow-y-auto px-4 py-3">
		{#if messages.length === 0}
			<p class="text-sm text-slate-400">
				Describe what you need — e.g. “a low-poly cartoon crown” or “free medieval props for my obby”.
			</p>
		{/if}
		{#each messages as msg}
			<div class={msg.role === 'user' ? 'text-right' : 'text-left'}>
				<span
					class={'inline-block max-w-[80%] rounded-2xl px-3 py-2 text-sm whitespace-pre-wrap ' +
						(msg.role === 'user' ? 'bg-indigo-600 text-white' : 'bg-slate-100 text-slate-900')}
				>
					{msg.content || (busy ? '…' : '')}
				</span>
				{#if msg.products?.length}
					<div class="mt-2 flex gap-2 overflow-x-auto pb-1">
						{#each msg.products as p (p.id)}
							<a
								href="/p/{p.slug}"
								class="w-36 shrink-0 rounded-lg border border-slate-200 bg-white p-2 text-left hover:border-indigo-400"
							>
								{#if p.thumb_url}
									<img src={p.thumb_url} alt={p.title} class="h-20 w-full rounded object-contain" />
								{/if}
								<p class="mt-1 line-clamp-2 text-xs font-medium">{p.title}</p>
								<p class="text-xs text-slate-500">{priceLabel(p)}</p>
							</a>
						{/each}
					</div>
				{/if}
			</div>
		{/each}
	</div>
	<form onsubmit={send} class="flex gap-2 border-t border-slate-200 p-3">
		<input
			bind:value={input}
			disabled={busy}
			placeholder="What are you looking for?"
			class="flex-1 rounded-full border border-slate-300 px-4 py-2 text-sm focus:border-indigo-500 focus:outline-none"
		/>
		<button
			type="submit"
			disabled={busy || !input.trim()}
			class="rounded-full bg-indigo-600 px-5 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-40"
		>
			Send
		</button>
	</form>
</div>
