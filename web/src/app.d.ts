import type { HTMLAttributes } from 'svelte/elements';

declare global {
	namespace App {}
}

// <model-viewer> web component (static/vendor/model-viewer.min.js)
declare module 'svelte/elements' {
	export interface SvelteHTMLElements {
		'model-viewer': HTMLAttributes<HTMLElement> & Record<string, unknown>;
	}
}

export {};
