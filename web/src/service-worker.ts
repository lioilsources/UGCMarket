/// <reference types="@sveltejs/kit" />
/// <reference lib="webworker" />

// Minimal PWA service worker: cache-first for the static build files,
// everything else (API, previews, downloads) always goes to the network.
import { build, files, version } from '$service-worker';

const sw = self as unknown as ServiceWorkerGlobalScope;
const CACHE = `ol1n-market-${version}`;
const ASSETS = [...build, ...files];

sw.addEventListener('install', (event) => {
	event.waitUntil(caches.open(CACHE).then((cache) => cache.addAll(ASSETS)));
});

sw.addEventListener('activate', (event) => {
	event.waitUntil(
		caches.keys().then((keys) => Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k))))
	);
});

sw.addEventListener('fetch', (event) => {
	if (event.request.method !== 'GET') return;
	const url = new URL(event.request.url);
	if (ASSETS.includes(url.pathname)) {
		event.respondWith(caches.match(event.request).then((hit) => hit ?? fetch(event.request)));
	}
});
