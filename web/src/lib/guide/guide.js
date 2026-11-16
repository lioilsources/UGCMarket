// Single source of the customer guide: the /guide page and the olin.now export both render this.
// Rule: a section is `verified: true` only after it was walked through in Roblox Studio on the Mac
// (see docs/STUDIO_RUNBOOK.md). Unverified sections render only in dev, never in production.
//
// Inline markup in strings: **bold** for exact Studio UI labels, `code` for file and property names.
// Blocks: { p }, { steps: [string | { text, img }] }, { table: { head, rows, rowIds? } },
//         { code, lang }, { img: { name, alt } }, { todo: [string] } (dev-only notes).
// Images are `web/static/guide/<name>.webp`, produced by bin/guide-shots; sizes live in shots.json.

/** @type {import('./types').Guide} */
export const GUIDE = {
	studioVersion: null, // fill in after the walkthrough, e.g. "0.738"
	verifiedOn: null, // e.g. "2026-11-17"
	sections: [
		{
			id: 'zip',
			title: "What's in the ZIP",
			verified: true, // describes our own package (api/internal/shop/download.go), not Studio
			blocks: [
				{ p: 'Every download is one folder named after the model. Inside you get:' },
				{
					table: {
						head: ['File', 'What it is for'],
						rows: [
							['`name.fbx`', 'The model for Roblox Studio. Up to 4,000 triangles, already sized in studs.'],
							['`name_tex.png`', 'The color texture for the FBX model (1024 × 1024 px).'],
							['`name.glb`', 'A more detailed version for other 3D tools. It is not made for Roblox limits.'],
							['`LICENSE.txt`', 'What you may do with the model, with your license key.'],
							['`README.txt`', 'A short summary and a link to this guide.']
						]
					}
				},
				{ p: 'For Roblox Studio you only need the `.fbx` and the `_tex.png` file.' }
			]
		},
		{
			id: 'before',
			title: 'Before you start',
			verified: false,
			blocks: [
				{
					todo: [
						'Studio version, sign-in, English UI.',
						'Does Import upload the mesh and texture to the signed-in account? Is it moderated? What does the user see?',
						'Warning: use an account you are comfortable uploading to.'
					]
				}
			]
		},
		{
			id: 'import',
			title: 'Import the model',
			verified: false,
			blocks: [
				{
					todo: [
						'Exact menu path to the 3D Importer in the current Studio.',
						'Settings the customer must change (if any). Is the model upright, facing forward, right size (bbox_studs)?',
						'What happens to the empty <Attachment>_Att object?'
					]
				}
			]
		},
		{
			id: 'texture',
			title: 'Apply the texture',
			verified: false,
			blocks: [
				{
					todo: [
						'Does the texture load by itself (trellis FBX references model_tex.png, sf3d FBX has none)?',
						'If not: shortest path to upload the PNG and set it on the MeshPart. How long does moderation take?'
					]
				}
			]
		},
		{
			id: 'accessory',
			title: 'Turn it into an accessory',
			verified: false,
			blocks: [
				{
					todo: [
						'Accessory Fitting Tool steps, exact labels. Does the model sit right without manual moves?',
						'Is Avatar Setup simpler for a beginner?',
						'Table category → accessory type (rows need ids cat-<category>, product pages link to them).'
					]
				}
			]
		},
		{
			id: 'try',
			title: 'Try it on a character',
			verified: false,
			blocks: [{ todo: ['Rig Builder dummy + drag the Accessory under it, or Test in Experience: which is shorter?'] }]
		},
		{
			id: 'game',
			title: 'Use it in your experience',
			verified: false,
			blocks: [{ todo: ['Minimal script that gives every player the accessory, checked in a playtest.'] }]
		},
		{
			id: 'troubleshooting',
			title: 'Troubleshooting',
			verified: false,
			blocks: [
				{ todo: ['No texture, wrong size, sits crooked, rejected by moderation (link /licenses for refunds).'] }
			]
		},
		{
			id: 'publishing',
			title: 'Publishing',
			verified: false,
			blocks: [{ todo: ['One paragraph + link to Roblox docs. No promises (R8).'] }]
		}
	]
};

// Product category → anchor on the guide page. Rows of the accessory table carry these ids.
/** @param {string | undefined} category */
export function guideAnchor(category) {
	return category ? `#cat-${category}` : '#import';
}
