import type { RequestLog } from '$lib/api/BeoApi';

// Searchable parts of a request log. Each one is shown as a "match in" chip
// and gets its own highlight list, so the UI can point at where a hit is.
export type LogField =
	| 'path'
	| 'method'
	| 'status'
	| 'query'
	| 'request_headers'
	| 'request_body'
	| 'response_headers'
	| 'response_body';

export const LOG_FIELD_LABELS: Record<LogField, string> = {
	path: 'URL',
	method: 'Method',
	status: 'Status',
	query: 'Query',
	request_headers: 'Req Header',
	request_body: 'Req Body',
	response_headers: 'Res Header',
	response_body: 'Res Body'
};

const REQUEST_FIELDS: LogField[] = ['path', 'method', 'query', 'request_headers', 'request_body'];
const RESPONSE_FIELDS: LogField[] = ['status', 'response_headers', 'response_body'];
const ALL_FIELDS: LogField[] = [...REQUEST_FIELDS, ...RESPONSE_FIELDS];
// Unprefixed terms skip the status code: searching `200` should look in the
// payload, not match every successful request. Use `status:200` for that.
const DEFAULT_FIELDS: LogField[] = ALL_FIELDS.filter((field) => field !== 'status');

// Prefixes accepted in `prefix:value` / `prefix=value`.
const FIELD_ALIASES: Record<string, LogField[]> = {
	path: ['path'],
	url: ['path', 'query'],
	method: ['method'],
	status: ['status'],
	query: ['query'],
	param: ['query'],
	params: ['query'],
	header: ['request_headers', 'response_headers'],
	headers: ['request_headers', 'response_headers'],
	body: ['request_body', 'response_body'],
	req: REQUEST_FIELDS,
	request: REQUEST_FIELDS,
	res: RESPONSE_FIELDS,
	response: RESPONSE_FIELDS,
	'req.header': ['request_headers'],
	'req.headers': ['request_headers'],
	'req.body': ['request_body'],
	'res.header': ['response_headers'],
	'res.headers': ['response_headers'],
	'res.body': ['response_body']
};

export const SEARCH_SYNTAX_HELP = [
	{ example: 'abc', description: 'match anywhere (URL, method, query, headers, body)' },
	{ example: 'body:abc  /  body=abc', description: 'request or response body' },
	{ example: 'req.body:abc  res.body:abc', description: 'only request / only response body' },
	{ example: 'header:content-type', description: 'request or response headers (req.header / res.header)' },
	{ example: 'query:brandId=wardah', description: 'query parameters (alias: param, params)' },
	{ example: 'path:users  method:post  status:500', description: 'URL path, HTTP method, status code' },
	{ example: 'req:abc  res:abc', description: 'anywhere in the request / the response' },
	{ example: '"two words"  body:"a b"', description: 'exact phrase containing spaces' },
	{ example: 'abc|def', description: 'either value (OR)' },
	{ example: '-abc  -status:200', description: 'exclude logs containing the value' }
];

export type SearchTerm = {
	// Alternatives joined by `|`, lowercased. A term matches when any of them does.
	values: string[];
	fields: LogField[];
	negate: boolean;
};

export type ParsedSearch = {
	terms: SearchTerm[];
};

export type LogMatch = {
	matched: boolean;
	fields: LogField[];
};

// Split on whitespace but keep quoted phrases together: `body:"a b" c` → [`body:"a b"`, `c`].
function tokenize(input: string): string[] {
	const tokens: string[] = [];
	let current = '';
	let inQuotes = false;
	for (const ch of input) {
		if (ch === '"') {
			inQuotes = !inQuotes;
			current += ch;
		} else if (/\s/.test(ch) && !inQuotes) {
			if (current) tokens.push(current);
			current = '';
		} else {
			current += ch;
		}
	}
	if (current) tokens.push(current);
	return tokens;
}

// Strip only wrapping quotes so JSON fragments like `"stock":5` stay searchable.
function unquote(value: string): string {
	if (value.length >= 2 && value.startsWith('"') && value.endsWith('"')) {
		return value.slice(1, -1);
	}
	return value;
}

export function parseSearch(input: string): ParsedSearch {
	const terms: SearchTerm[] = [];

	for (let token of tokenize(input.trim())) {
		let negate = false;
		if (token.startsWith('-') && token.length > 1) {
			negate = true;
			token = token.slice(1);
		}

		let fields = DEFAULT_FIELDS;
		let value = token;
		// `prefix:value` or `prefix=value`, only when the prefix is a known field
		// so plain text such as `brandId=wardah` still searches everywhere.
		const sep = token.search(/[:=]/);
		if (sep > 0 && !token.slice(0, sep).includes('"')) {
			const prefix = token.slice(0, sep).toLowerCase();
			if (FIELD_ALIASES[prefix]) {
				fields = FIELD_ALIASES[prefix];
				value = token.slice(sep + 1);
			}
		}

		const values = unquote(value)
			.toLowerCase()
			.split('|')
			.filter((v) => v !== '');
		if (values.length === 0) continue;

		terms.push({ values, fields, negate });
	}

	return { terms };
}

function headersToText(raw: string): string {
	if (!raw) return '';
	try {
		const parsed = JSON.parse(raw);
		if (parsed && typeof parsed === 'object') {
			return Object.entries(parsed)
				.map(([key, value]) => `${key}: ${Array.isArray(value) ? value.join(', ') : String(value)}`)
				.join('\n');
		}
	} catch {
		// not JSON, search the raw string
	}
	return raw;
}

// Raw body plus the pretty-printed JSON the panel shows, so a phrase copied
// from either form matches.
function bodyToText(raw: string): string {
	if (!raw) return '';
	try {
		return `${raw}\n${JSON.stringify(JSON.parse(raw), null, 2)}`;
	} catch {
		return raw;
	}
}

const textCache = new WeakMap<RequestLog, Record<LogField, string>>();

function fieldTexts(log: RequestLog): Record<LogField, string> {
	const cached = textCache.get(log);
	if (cached) return cached;

	const texts: Record<LogField, string> = {
		path: (log.path || '').toLowerCase(),
		method: (log.method || '').toLowerCase(),
		status: String(log.response_status ?? ''),
		query: (log.query_params || '').toLowerCase(),
		request_headers: headersToText(log.request_headers).toLowerCase(),
		request_body: bodyToText(log.request_body).toLowerCase(),
		response_headers: headersToText(log.response_headers).toLowerCase(),
		response_body: bodyToText(log.response_body).toLowerCase()
	};
	textCache.set(log, texts);
	return texts;
}

// All positive terms must hit at least one of their fields, and no negative
// term may hit any of its fields. `fields` lists every place a positive term hit.
export function matchLog(log: RequestLog, search: ParsedSearch): LogMatch {
	if (search.terms.length === 0) return { matched: true, fields: [] };

	const texts = fieldTexts(log);
	const hitFields = new Set<LogField>();

	for (const term of search.terms) {
		const hits = term.fields.filter((field) =>
			term.values.some((value) => texts[field].includes(value))
		);
		if (term.negate) {
			if (hits.length > 0) return { matched: false, fields: [] };
		} else {
			if (hits.length === 0) return { matched: false, fields: [] };
			hits.forEach((field) => hitFields.add(field));
		}
	}

	return { matched: true, fields: ALL_FIELDS.filter((field) => hitFields.has(field)) };
}

// Values to highlight inside one field of the expanded log.
export function highlightTermsFor(search: ParsedSearch | null | undefined, field: LogField): string[] {
	if (!search) return [];
	return search.terms
		.filter((term) => !term.negate && term.fields.includes(field))
		.flatMap((term) => term.values);
}

export function isResponseOnlyMatch(fields: LogField[]): boolean {
	return fields.length > 0 && fields.every((field) => RESPONSE_FIELDS.includes(field));
}

export type HighlightSegment = { text: string; match: boolean };

// Split text into plain / matched segments (case-insensitive, overlapping hits merged).
export function highlightSegments(text: string, terms: string[]): HighlightSegment[] {
	if (!text || terms.length === 0) return [{ text: text ?? '', match: false }];

	const lower = text.toLowerCase();
	const ranges: Array<[number, number]> = [];
	for (const term of terms) {
		if (!term) continue;
		let from = 0;
		let index: number;
		while ((index = lower.indexOf(term, from)) !== -1) {
			ranges.push([index, index + term.length]);
			from = index + term.length;
		}
	}
	if (ranges.length === 0) return [{ text, match: false }];

	ranges.sort((a, b) => a[0] - b[0]);
	const merged: Array<[number, number]> = [];
	for (const range of ranges) {
		const last = merged[merged.length - 1];
		if (last && range[0] <= last[1]) {
			last[1] = Math.max(last[1], range[1]);
		} else {
			merged.push([...range]);
		}
	}

	const segments: HighlightSegment[] = [];
	let cursor = 0;
	for (const [start, end] of merged) {
		if (start > cursor) segments.push({ text: text.slice(cursor, start), match: false });
		segments.push({ text: text.slice(start, end), match: true });
		cursor = end;
	}
	if (cursor < text.length) segments.push({ text: text.slice(cursor), match: false });
	return segments;
}
