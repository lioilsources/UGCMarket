import { env } from '$env/dynamic/private';
import type { RequestHandler } from './$types';

const API_URL = (env.API_URL ?? 'http://localhost:8096').replace(/\/$/, '');

// Headers that market-api needs (session cookie, country and IP from Cloudflare).
const FORWARD = ['cookie', 'content-type', 'accept', 'range', 'cf-ipcountry', 'cf-connecting-ip', 'x-forwarded-for'];

// The whole /api/* goes 1:1 to market-api — the web has no direct DB access.
const proxy: RequestHandler = async ({ request, params, url }) => {
	const headers = new Headers();
	for (const name of FORWARD) {
		const v = request.headers.get(name);
		if (v) headers.set(name, v);
	}
	const hasBody = request.method !== 'GET' && request.method !== 'HEAD';
	const upstream = await fetch(`${API_URL}/api/${params.path}${url.search}`, {
		method: request.method,
		headers,
		body: hasBody ? request.body : undefined,
		// @ts-expect-error — Node fetch requires duplex for a streamed body
		duplex: hasBody ? 'half' : undefined,
		redirect: 'manual'
	});
	const out = new Headers();
	for (const [k, v] of upstream.headers) {
		if (!['content-encoding', 'content-length', 'transfer-encoding', 'connection'].includes(k)) out.append(k, v);
	}
	for (const c of upstream.headers.getSetCookie()) out.append('set-cookie', c);
	return new Response(upstream.body, { status: upstream.status, headers: out });
};

export const GET = proxy;
export const POST = proxy;
export const PATCH = proxy;
export const PUT = proxy;
export const DELETE = proxy;
