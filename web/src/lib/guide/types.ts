export type GuideImage = { name: string; alt: string };

export type GuideStep = string | { text: string; img?: GuideImage };

export type GuideBlock =
	| { p: string }
	| { steps: GuideStep[] }
	| { table: { head: string[]; rows: string[][]; rowIds?: string[] } }
	| { code: string; lang?: string }
	| { img: GuideImage }
	| { todo: string[] };

export type GuideSection = { id: string; title: string; verified: boolean; blocks: GuideBlock[] };

export type Guide = { studioVersion: string | null; verifiedOn: string | null; sections: GuideSection[] };
