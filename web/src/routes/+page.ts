import type { FilterSchema, Product } from '$lib/types';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const [facets, products] = await Promise.all([
		fetch('/api/catalog/facets').then((r) => (r.ok ? (r.json() as Promise<FilterSchema[]>) : [])),
		fetch('/api/products').then((r) => (r.ok ? (r.json() as Promise<Product[]>) : []))
	]);
	return { facets, products };
};
