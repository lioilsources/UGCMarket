import { browser } from '$app/environment';
import type { Product } from './types';

// Digital goods: each asset in the cart at most once, no quantity.
export type CartItem = Pick<Product, 'id' | 'slug' | 'title' | 'price_cents' | 'currency' | 'tier' | 'thumb_url'>;

const STORAGE_KEY = 'ol1n-market-cart';

function load(): CartItem[] {
	if (!browser) return [];
	try {
		return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '[]');
	} catch {
		return [];
	}
}

class Cart {
	items = $state<CartItem[]>(load());

	#save() {
		if (browser) localStorage.setItem(STORAGE_KEY, JSON.stringify(this.items));
	}

	has(id: number) {
		return this.items.some((i) => i.id === id);
	}

	add(p: Product | CartItem) {
		if (this.has(p.id)) return;
		const { id, slug, title, price_cents, currency, tier, thumb_url } = p;
		this.items.push({ id, slug, title, price_cents, currency, tier, thumb_url });
		this.#save();
	}

	remove(id: number) {
		this.items = this.items.filter((i) => i.id !== id);
		this.#save();
	}

	clear() {
		this.items = [];
		this.#save();
	}

	get count() {
		return this.items.length;
	}

	get totalCents() {
		return this.items.reduce((sum, i) => sum + (i.tier === 'free' ? 0 : i.price_cents), 0);
	}
}

export const cart = new Cart();
