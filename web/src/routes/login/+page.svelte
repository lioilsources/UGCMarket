<script lang="ts">
	import { page } from '$app/state';
	import { apiError } from '$lib/format';

	let email = $state('');
	let busy = $state(false);
	let sent = $state(false);
	let message = $state(page.url.searchParams.get('error') === 'link' ? 'That sign-in link is invalid or expired. Request a new one.' : '');

	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		message = '';
		const res = await fetch('/api/auth/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email, next: page.url.searchParams.get('next') ?? '/me' })
		});
		busy = false;
		if (res.ok) sent = true;
		else message = await apiError(res);
	}
</script>

<svelte:head><title>Sign in — Ol1n Market</title></svelte:head>

<div class="mx-auto w-full max-w-sm space-y-4 p-4 md:p-8">
	<h1 class="text-2xl font-bold">Sign in</h1>
	{#if sent}
		<p class="rounded-lg bg-emerald-50 px-3 py-2 text-emerald-800">
			Check <strong>{email}</strong> — we sent you a sign-in link. It works once and expires in 30 minutes.
		</p>
	{:else}
		<p class="text-sm text-slate-600">No password. We'll e-mail you a one-time sign-in link; your purchases live under that address.</p>
		<form onsubmit={submit} class="space-y-3">
			<input
				type="email"
				required
				bind:value={email}
				placeholder="you@example.com"
				class="w-full rounded-lg border border-slate-300 px-3 py-2 focus:border-indigo-500 focus:outline-none"
			/>
			<button
				disabled={busy}
				class="w-full rounded-lg bg-indigo-600 py-2 font-medium text-white hover:bg-indigo-700 disabled:opacity-40"
			>
				Send sign-in link
			</button>
		</form>
	{/if}
	{#if message}
		<p class="text-sm text-red-600">{message}</p>
	{/if}
</div>
