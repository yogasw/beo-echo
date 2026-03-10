import type { Replay } from '$lib/types/Replay';
import { jsonParser } from './jsonLogs';
import { curlParser } from './curl';

export interface ParserResult {
	parsed: Partial<Replay>;
	importType: string;
	displayName?: string;
	rawText: string;
}

export interface Parser {
	name: string;
	displayName: string;
	canParse: (text: string) => boolean;
	parse: (text: string) => ParserResult;
}

export const parsers: Parser[] = [
	{
		name: 'curl',
		displayName: 'cURL',
		canParse: (text: string) => text.trim().toLowerCase().startsWith('curl'),
		parse: curlParser.parse
	},
	{
		name: 'json-logs',
		displayName: 'JSON Logs',
		canParse: (text: string) => {
			const trimmed = text.trim();
			return trimmed.startsWith('{') || trimmed.startsWith('[');
		},
		parse: jsonParser.parse
	}
];

export function parseImportText(text: string): ParserResult {
	const trimmed = text.trim();
	
	for (const parser of parsers) {
		if (parser.canParse(trimmed)) {
			const result = parser.parse(trimmed);
			return { ...result, displayName: parser.displayName };
		}
	}

	throw new Error("Format not supported. Please paste a valid JSON or cURL command.");
}

export * from './jsonLogs';
export * from './curl';
