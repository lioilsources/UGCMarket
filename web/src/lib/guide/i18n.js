// Translations of the guide (olin.now ships it in every site language; market.ol1n.com/guide is English).
// guide.js stays the only source of structure; a translation is a flat map key -> string with the
// same keys as guideStrings(GUIDE). Code blocks, image names and anchors are never translated.
//
// Keys: <section>.title, <section>.<block>.p, .steps.<i>, .alt, .steps.<i>.alt, .th.<j>, .td.<r>.<c>

/** @param {import('./types').Guide} guide @returns {Record<string, string>} */
export function guideStrings(guide) {
	/** @type {Record<string, string>} */
	const out = {};
	for (const s of guide.sections) {
		if (!s.verified) continue;
		out[`${s.id}.title`] = s.title;
		s.blocks.forEach((b, bi) => {
			const k = `${s.id}.${bi}`;
			if ('p' in b) out[`${k}.p`] = b.p;
			if ('steps' in b)
				b.steps.forEach((st, i) => {
					if (typeof st === 'string') out[`${k}.steps.${i}`] = st;
					else {
						out[`${k}.steps.${i}`] = st.text;
						if (st.img) out[`${k}.steps.${i}.alt`] = st.img.alt;
					}
				});
			if ('img' in b) out[`${k}.alt`] = b.img.alt;
			if ('table' in b) {
				b.table.head.forEach((h, j) => (out[`${k}.th.${j}`] = h));
				b.table.rows.forEach((r, ri) => r.forEach((c, ci) => (out[`${k}.td.${ri}.${ci}`] = c)));
			}
		});
	}
	return out;
}

/**
 * A copy of the guide with every string replaced from `tr`. Throws when a key is missing,
 * so a half-done translation never ships.
 * @param {import('./types').Guide} guide
 * @param {Record<string, string>} tr
 * @returns {import('./types').Guide}
 */
export function applyStrings(guide, tr) {
	/** @param {string} key */
	const t = (key) => {
		if (typeof tr[key] !== 'string') throw new Error(`missing translation: ${key}`);
		return tr[key];
	};
	return {
		...guide,
		sections: guide.sections
			.filter((s) => s.verified)
			.map((s) => ({
				...s,
				title: t(`${s.id}.title`),
				blocks: s.blocks.map((b, bi) => {
					const k = `${s.id}.${bi}`;
					if ('p' in b) return { p: t(`${k}.p`) };
					if ('steps' in b)
						return {
							steps: b.steps.map((st, i) =>
								typeof st === 'string'
									? t(`${k}.steps.${i}`)
									: { text: t(`${k}.steps.${i}`), img: st.img && { ...st.img, alt: t(`${k}.steps.${i}.alt`) } }
							)
						};
					if ('img' in b) return { img: { ...b.img, alt: t(`${k}.alt`) } };
					if ('table' in b)
						return {
							table: {
								...b.table,
								head: b.table.head.map((_, j) => t(`${k}.th.${j}`)),
								rows: b.table.rows.map((r, ri) => r.map((_, ci) => t(`${k}.td.${ri}.${ci}`)))
							}
						};
					return b;
				})
			}))
	};
}
