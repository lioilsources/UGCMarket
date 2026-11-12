import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const API_URL = (env.API_URL ?? 'http://localhost:8096').replace(/\/$/, '');

/** Target of the e-mail link: verifies the token in market-api, passes on the session cookie and redirects. */
export const GET: RequestHandler = async ({ url, cookies }) => {
	const token = url.searchParams.get('token') ?? '';
	const res = await fetch(`${API_URL}/api/auth/verify?token=${encodeURIComponent(token)}`);
	if (!res.ok) redirect(303, '/login?error=link');

	const match = /om_session=([^;]+)/.exec(res.headers.getSetCookie().join('\n'));
	if (!match) redirect(303, '/login?error=link');
	cookies.set('om_session', match[1], {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: url.protocol === 'https:',
		maxAge: 60 * 60 * 24 * 60
	});

	const next = url.searchParams.get('next') ?? '/me';
	redirect(303, next.startsWith('/') && !next.startsWith('//') ? next : '/me');
};
