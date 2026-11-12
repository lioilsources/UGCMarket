export type FilterKind = 'set' | 'range' | 'bool';

export interface FilterSchema {
	attr_key: string;
	label: string;
	kind: FilterKind;
	config: { values?: string[]; min?: number; max?: number; unit?: string };
}

/** Active filter values: set -> string[], range -> {min,max}, bool -> boolean */
export type ActiveFilters = Record<string, string[] | { min?: number; max?: number } | boolean>;

export interface AssetAttributes {
	category?: string;
	style?: string;
	collection?: string;
	tags?: string[];
	formats?: string[];
	polycount?: number;
	texture_px?: number;
	watertight?: boolean;
	roblox_attachment?: string;
	bbox_studs?: number[];
	[key: string]: unknown;
}

/** Public shape of a product from market-api (catalog.Card). */
export interface Product {
	id: number;
	slug: string;
	title: string;
	description: string;
	price_cents: number;
	currency: string;
	tier: 'free' | 'paid';
	status: 'draft' | 'listed' | 'archived';
	attributes: AssetAttributes;
	thumb_url?: string;
	viewer_url?: string;
	created_at: string;
}

export interface ProductDetail extends Product {
	files: { format: string; bytes: number }[];
	owned: boolean;
}

export interface Me {
	id: number;
	email: string;
	is_admin: boolean;
}

export interface License {
	id: number;
	product_id: number;
	license_type: 'personal' | 'commercial';
	license_key: string;
	issued_at: string;
	title: string;
	slug: string;
	tier: string;
}

export interface ChatMsg {
	role: 'user' | 'assistant';
	content: string;
	products?: Product[];
}
