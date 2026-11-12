import { redirect } from '@sveltejs/kit';
import type { License } from '$lib/types';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, parent }) => {
	const { me } = await parent();
	if (!me) redirect(303, '/login?next=/me');
	const res = await fetch('/api/me/licenses');
	return { licenses: res.ok ? ((await res.json()) as License[]) : [] };
};
