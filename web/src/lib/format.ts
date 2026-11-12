import type { Product } from './types';

export function priceLabel(p: Pick<Product, 'tier' | 'price_cents' | 'currency'>): string {
	if (p.tier === 'free') return 'Free';
	return new Intl.NumberFormat('en-US', { style: 'currency', currency: p.currency }).format(
		p.price_cents / 100
	);
}

export function bytesLabel(n: number): string {
	if (n < 1024) return `${n} B`;
	if (n < 1024 * 1024) return `${(n / 1024).toFixed(0)} KB`;
	return `${(n / 1024 / 1024).toFixed(1)} MB`;
}

/** Error message from market-api ({error}) or the status. */
export async function apiError(res: Response): Promise<string> {
	try {
		const body = await res.json();
		if (body?.error) return body.error;
	} catch {
		// non-JSON response
	}
	return `Request failed (${res.status})`;
}
