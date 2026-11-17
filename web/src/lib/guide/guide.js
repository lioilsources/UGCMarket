// Single source of the customer guide: the /guide page and the olin.now export both render this.
// Rule: a section is `verified: true` only after it was walked through in Roblox Studio on the Mac
// (see docs/STUDIO_RUNBOOK.md). Unverified sections render only in dev, never in production.
//
// Inline markup in strings: **bold** for exact Studio UI labels, `code` for file and property names,
// [text](url) for links (absolute URLs or #anchors, so the olin.now export works too).
// Blocks: { p }, { steps: [string | { text, img }] }, { table: { head, rows, rowIds? } },
//         { code, lang }, { img: { name, alt } }, { todo: [string] } (dev-only notes).
// Images are `web/static/guide/<name>.webp`, produced by bin/guide-shots; sizes live in shots.json.

/** @type {import('./types').Guide} */
export const GUIDE = {
	studioVersion: '0.747.0.7471481',
	verifiedOn: '2026-11-17',
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
							['`name.fbx`', 'The model for Roblox Studio. Up to 4,000 triangles. Most models are already sized in studs; see Troubleshooting if yours is not.'],
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
			verified: true,
			blocks: [
				{
					steps: [
						'Install Roblox Studio and sign in. Buttons may move a little in newer versions than the one named above.',
						'Open a place to work in. A new **Baseplate** place is enough.',
						'Unzip the download. The ZIP opens into a folder with the model name, and the files are in another folder with the same name inside it.'
					]
				},
				{
					p: 'Importing uploads the model and the texture to the Roblox account you are signed in with (**Upload to Roblox** is on by default). Use an account you are happy to upload to.'
				}
			]
		},
		{
			id: 'import',
			title: 'Import the model',
			verified: true,
			blocks: [
				{
					steps: [
						{
							text: 'On the **Home** tab, click **Import**. Do not use the **Import** button inside the Terrain Editor; that one is for height maps.',
							img: { name: '01-import-button', alt: 'Import button on the Home tab' }
						},
						'In the file dialog, pick the `.fbx` file from the ZIP.',
						{
							text: 'The **Import Preview** window opens. The default settings are fine. Before you click **Import**, fix the texture as shown in the next section.',
							img: { name: '02-import-preview', alt: 'Import Preview window with default settings' }
						}
					]
				},
				{
					p: 'After the import, the **Import Queue** panel shows a green check, and a new `Model` with a `MeshPart` called `Handle` appears in `Workspace`. The model stands upright and faces the front.'
				}
			]
		},
		{
			id: 'texture',
			title: 'Apply the texture',
			verified: true,
			blocks: [
				{
					p: 'Most models show a `base_color_texture` row in **Import Preview** with a red mark: “Could not read the color map texture file associated with this object.” You can fix it right there.'
				},
				{
					steps: [
						'In the tree on the left, click `base_color_texture` (under `world` ⟩ `Handle`).',
						{
							text: 'Next to **Color File Path**, click the folder icon.',
							img: { name: '03-texture-error', alt: 'Texture row with the error and the folder icon next to Color File Path' }
						},
						{
							text: 'Pick the `_tex.png` file from the ZIP. The row gets a check mark. Click **Import**.',
							img: { name: '04-texture-picked', alt: 'Texture row checked after picking the PNG file' }
						}
					]
				},
				{
					p: 'Some models have no `base_color_texture` row at all and come in white. Add the texture after the import:'
				},
				{
					steps: [
						'On the **Home** tab, click **Assets** to open the **Asset Manager**.',
						'Click **Import** in the Asset Manager and pick the `_tex.png` file. In the **Import Queue**, click **Start Import**.',
						{
							text: 'When the green check appears, find the image in the Asset Manager and note its number in the **ID** column.',
							img: { name: '11-asset-manager-id', alt: 'Asset Manager with the image ID' }
						},
						'In **Explorer**, select the `Handle` part of the model. In **Properties**, type `texture` in the filter.',
						'Click **TextureContent**, type `rbxassetid://` followed by the ID (for example `rbxassetid://1234567890`) and press Enter.'
					]
				}
			]
		},
		{
			id: 'accessory',
			title: 'Turn it into an accessory',
			verified: true,
			blocks: [
				{
					steps: [
						{
							text: 'On the **Avatar** tab, click **Accessory**. The **Accessory Fitting Tool** panel opens.',
							img: { name: '05-avatar-accessory', alt: 'Accessory button on the Avatar tab' }
						},
						{
							text: 'Click the **Part** field, then click the model in the viewport. The field shows `Handle`. Click **Continue**.',
							img: { name: '06-aft-select', alt: 'Handle selected in the Accessory Fitting Tool' }
						},
						{
							text: 'Choose **Accessory**, then the accessory type from the table below. Leave **Select Body Type** as it is. A mannequin appears and the model moves to its place on it. Click **Continue**.',
							img: { name: '07-aft-type', alt: 'Accessory type Back chosen, the model sits on the mannequin' }
						},
						{
							text: 'The blue box shows the space the accessory may take. Click **Generate MeshPart Accessory** at the bottom of the panel.',
							img: { name: '08-aft-generate', alt: 'Generate MeshPart Accessory button' }
						},
						'A small window says the accessory was generated into Workspace. It may open behind the main Studio window. Click **OK**.'
					]
				},
				{
					p: 'The new `Accessory` is in `Workspace` and is named `Handle`, after the part. Rename it (for example `KabutoHelmet`) and delete the original imported `Model`, which has almost the same name, so you do not mix them up.'
				},
				{
					table: {
						head: ['Model category', 'Accessory type', 'Where it sits'],
						rows: [
							['Hat', '**Hat**', 'On the head'],
							['Helmet', '**Hat**', 'On the head'],
							['Back', '**Back**', 'On the back']
						],
						rowIds: ['cat-hat', 'cat-helmet', 'cat-back']
					}
				}
			]
		},
		{
			id: 'try',
			title: 'Try it on a character',
			verified: true,
			blocks: [
				{
					steps: [
						{
							text: 'On the **Avatar** tab, click **Character**. In **Generate Rig**, keep **R15** and click **Block Avatar (2012)**. A figure called `Rig` appears in `Workspace`.',
							img: { name: '09-rig-builder', alt: 'Generate Rig with Block Avatar (2012)' }
						},
						{
							text: 'In **Explorer**, drag your accessory onto `Rig`. It snaps to its place without any moving by hand.',
							img: { name: '10-helmet-on-rig', alt: 'Helmet on the test figure' }
						}
					]
				},
				{ img: { name: '12-backpack-on-rig', alt: 'Backpack on the back of the test figure' } }
			]
		},
		{
			id: 'game',
			title: 'Use it in your experience',
			verified: true,
			blocks: [
				{
					steps: [
						'Move the accessory to `ServerStorage` and give it a name without spaces, for example `KabutoHelmet`.',
						'In `ServerScriptService`, add a `Script` with the code below. Change the name in `WaitForChild` to yours.',
						'Click **Play**. Your character wears the accessory.'
					]
				},
				{
					code: `local Players = game:GetService("Players")
local ServerStorage = game:GetService("ServerStorage")

local accessory = ServerStorage:WaitForChild("KabutoHelmet")

Players.PlayerAdded:Connect(function(player)
	player.CharacterAppearanceLoaded:Connect(function(character)
		local humanoid = character:WaitForChild("Humanoid")
		humanoid:AddAccessory(accessory:Clone())
	end)
end)`,
					lang: 'lua'
				}
			]
		},
		{
			id: 'troubleshooting',
			title: 'Troubleshooting',
			verified: true,
			blocks: [
				{
					table: {
						head: ['Problem', 'What to do'],
						rows: [
							['The model is white or grey.', 'The texture was not applied. See [Apply the texture](#texture).'],
							[
								'The model is much bigger than the character.',
								'Select the imported `Model` in **Explorer**, type `scale` in the **Properties** filter and set **Scale** to a smaller number (for example `0.3`). Do this before the Accessory Fitting Tool.'
							],
							['**Import** opened a height map dialog.', 'You used the Terrain Editor. Close it and use **Home** ⟩ **Import**.'],
							['The fitting tool does not pick the model.', 'Click the **Part** field first, then the model in the viewport.'],
							['The “generated” window does not show up.', 'It is behind the main Studio window. Move the main window or switch to it with the keyboard.']
						]
					}
				},
				{ p: 'If a model cannot be used at all, see [licenses and refunds](https://market.ol1n.com/licenses).' }
			]
		},
		{
			id: 'publishing',
			title: 'Publishing',
			verified: true, // only a link to the Roblox docs, nothing we claim about Studio (R8)
			blocks: [
				{
					p: 'Selling or sharing an accessory on the Marketplace has its own rules, checks and fees, set by Roblox. Read [Publish to Marketplace](https://create.roblox.com/docs/marketplace/publishing-to-marketplace) in the Roblox documentation before you try. Check also what your license allows.'
				}
			]
		}
	]
};

// Product category → anchor on the guide page. Rows of the accessory table carry these ids.
/** @param {string | undefined} category */
export function guideAnchor(category) {
	return category ? `#cat-${category}` : '#import';
}
