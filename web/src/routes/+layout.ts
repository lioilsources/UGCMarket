import type { Me } from '$lib/types';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ fetch }) => {
	const res = await fetch('/api/me');
	const me: Me | null = res.ok ? await res.json() : null;
	return { me };
};
