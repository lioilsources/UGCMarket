import { error } from '@sveltejs/kit';
import type { ProductDetail } from '$lib/types';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	const res = await fetch(`/api/products/${encodeURIComponent(params.slug)}`);
	if (res.status === 404) error(404, 'Asset not found');
	if (!res.ok) error(res.status, 'Could not load the asset');
	return { product: (await res.json()) as ProductDetail };
};
