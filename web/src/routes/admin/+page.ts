import { error, redirect } from '@sveltejs/kit';
import type { Product } from '$lib/types';
import type { PageLoad } from './$types';

export interface AdminOrder {
	id: number;
	email: string;
	total_cents: number;
	currency: string;
	status: string;
	provider: string;
	tax_country: string | null;
	item_count: number;
	created_at: string;
}

export const load: PageLoad = async ({ fetch, parent }) => {
	const { me } = await parent();
	if (!me) redirect(303, '/login?next=/admin');
	if (!me.is_admin) error(403, 'Admins only');
	const [products, ingest, orders] = await Promise.all([
		fetch('/api/admin/products').then((r) => r.json() as Promise<Product[]>),
		fetch('/api/admin/ingest').then((r) => r.json()),
		fetch('/api/admin/orders').then((r) => r.json() as Promise<{ orders: AdminOrder[]; downloads: number }>)
	]);
	return { products, ingest, orders };
};
