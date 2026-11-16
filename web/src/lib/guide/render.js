// Renders the guide (guide.js) to an HTML string. Plain JS so bin/ scripts can import it with node.
import shots from './shots.json' with { type: 'json' };

/** @param {string} s */
const esc = (s) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');

/** Escape, then apply **bold** and `code`. @param {string} s */
export function inline(s) {
	return esc(s)
		.replace(/`([^`]+)`/g, '<code>$1</code>')
		.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
}

/**
 * @param {import('./types').GuideImage} img
 * @param {string} base
 */
function image(img, base) {
	/** @type {Record<string, {w: number, h: number}>} */
	const sizes = shots;
	const s = sizes[img.name];
	const dims = s ? ` width="${s.w}" height="${s.h}"` : '';
	return `<img src="${base}${esc(img.name)}.webp" alt="${esc(img.alt)}"${dims} loading="lazy" data-zoom>`;
}

/**
 * @param {import('./types').GuideBlock} b
 * @param {{ imgBase: string, drafts: boolean }} o
 */
function block(b, o) {
	if ('p' in b) return `<p>${inline(b.p)}</p>`;
	if ('steps' in b) {
		const li = b.steps.map((st) =>
			typeof st === 'string'
				? `<li>${inline(st)}</li>`
				: `<li>${inline(st.text)}${st.img ? image(st.img, o.imgBase) : ''}</li>`
		);
		return `<ol class="steps">${li.join('')}</ol>`;
	}
	if ('table' in b) {
		const { head, rows, rowIds } = b.table;
		const th = head.map((h) => `<th>${inline(h)}</th>`).join('');
		const tr = rows.map((r, i) => {
			const id = rowIds?.[i] ? ` id="${esc(rowIds[i])}"` : '';
			return `<tr${id}>${r.map((c) => `<td>${inline(c)}</td>`).join('')}</tr>`;
		});
		return `<div class="table"><table><thead><tr>${th}</tr></thead><tbody>${tr.join('')}</tbody></table></div>`;
	}
	if ('code' in b) return `<pre><code>${esc(b.code)}</code></pre>`;
	if ('img' in b) return `<figure>${image(b.img, o.imgBase)}<figcaption>${inline(b.img.alt)}</figcaption></figure>`;
	if ('todo' in b) {
		if (!o.drafts) return '';
		return `<ul class="todo">${b.todo.map((t) => `<li>TODO: ${inline(t)}</li>`).join('')}</ul>`;
	}
	return '';
}

/**
 * @param {import('./types').Guide} guide
 * @param {{ imgBase?: string, drafts?: boolean }} [opts]
 *   imgBase: URL prefix of the screenshots; drafts: also render unverified sections (dev only).
 */
export function renderGuide(guide, opts = {}) {
	const o = { imgBase: opts.imgBase ?? '/guide/', drafts: opts.drafts ?? false };
	const sections = guide.sections.filter((s) => s.verified || o.drafts);
	const toc = sections.map((s) => `<li><a href="#${s.id}">${inline(s.title)}</a></li>`).join('');
	const body = sections.map((s) => {
		const badge = s.verified ? '' : ' <span class="draft">draft</span>';
		return `<section id="${s.id}"><h2>${inline(s.title)}${badge}</h2>${s.blocks.map((b) => block(b, o)).join('')}</section>`;
	});
	return { toc: `<ol class="toc">${toc}</ol>`, body: body.join(''), complete: sections.length === guide.sections.length };
}
