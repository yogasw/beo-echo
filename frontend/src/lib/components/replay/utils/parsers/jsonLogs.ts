import type { Replay } from '$lib/types/Replay';
import type { BodyTypeHttp,ReplayMetadata } from '$lib/types/Replay';

// --- Shared Types ---

export interface Header {
	key: string;
	value: string;
}

export interface ImportConfig {
	url: string;
	method: string;
	headers: string;
	body: string;
}

export interface ImportResult {
	url?: string;
	method?: string;
	headers?: Header[];
	body?: string;
	metadata?: ReplayMetadata;
	config: ImportConfig;
	availableKeys: string[];
}

// --- Utils ---

export const METHODS = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'];

/** Infer HTTP body type from body content. Extend here to support more types. */
export function inferBodyType(body?: string): BodyTypeHttp {
	return body ? 'raw' : 'none';
}

export function extractKeyPaths(data: any, prefix = ''): string[] {
	if (!data || typeof data !== 'object') return [];
	let paths: string[] = [];
	if (Array.isArray(data)) {
		for (let i = 0; i < data.length; i++) {
			paths.push(`${prefix}[${i}]`);
			paths = paths.concat(extractKeyPaths(data[i], `${prefix}[${i}]`));
		}
	} else {
		for (const [key, val] of Object.entries(data)) {
			const newPrefix = prefix ? `${prefix}.${key}` : key;
			paths.push(newPrefix);
			paths = paths.concat(extractKeyPaths(val, newPrefix));
		}
	}
	return paths;
}

export function getNestedValue(data: any, path: string): any {
	if (!path || !data) return undefined;
	const standardizedPath = path.replace(/\[(\w+)\]/g, '.$1');
	const parts = standardizedPath.split('.');
	let current = data;
	for (const part of parts) {
		if (current === null || current === undefined) return undefined;
		if (part === '') continue;
		current = current[part];
	}
	return current;
}

export function isValidUrl(value: any): boolean {
	if (typeof value !== 'string') return false;
	try {
		const url = new URL(value);
		return ['http:', 'https:'].includes(url.protocol);
	} catch {
		return value.startsWith('http://') || value.startsWith('https://');
	}
}

export function isValidMethod(value: any): boolean {
	return typeof value === 'string' && METHODS.includes(value.toUpperCase());
}

export function isValidHeaders(value: any): boolean {
	if (!value || typeof value !== 'object') return false;
	if (Array.isArray(value)) {
		if (value.length === 0) return false;
		const item = value[0];
		return item && typeof item === 'object' && (('key' in item) || ('name' in item) || Object.keys(item).length === 1);
	}
	const values = Object.values(value);
	if (values.length === 0) return false;
	const hasComplexValue = values.some(v => v !== null && typeof v === 'object');
	if (hasComplexValue) return false;
	return true;
}

export function resolveValue(data: any, configStr: string): any {
	if (!configStr) return undefined;
	const directValue = getNestedValue(data, configStr);
	if (directValue !== undefined) return directValue;
	if (configStr.includes('$')) {
		return configStr.replace(/(\$([a-zA-Z0-9_.[\]]+))|(\$\{([a-zA-Z0-9_.[\]]+)\})/g, (match, p1, p2, p3, p4) => {
			const key = p2 || p4;
			const val = getNestedValue(data, key);
			if (val !== undefined) return String(val);
			return match;
		});
	}
	return undefined;
}

export function findBestMatch(availableKeys: string[], keywords: string[]): string {
	if (!keywords) return '';
	const lowerKeys = availableKeys.map(k => k.toLowerCase());
	for (const keyword of keywords) {
		if (availableKeys.includes(keyword)) return keyword;
		const index = lowerKeys.indexOf(keyword.toLowerCase());
		if (index !== -1) return availableKeys[index];
		const suffixMatch = availableKeys.find(k => k.toLowerCase().endsWith(`.${keyword.toLowerCase()}`));
		if (suffixMatch) return suffixMatch;
	}
	return '';
}

function findValidatedMatch(
	data: any,
	availableKeys: string[],
	defaultKeywords: string[],
	validator: (val: any) => boolean
): string {
	if (!defaultKeywords) return '';
	for (const keyword of defaultKeywords) {
		if (keyword.includes('$')) {
			const val = resolveValue(data, keyword);
			if (val !== undefined && validator(val)) return keyword;
		} else {
			const candidateKey = findBestMatch(availableKeys, [keyword]);
			if (candidateKey) {
				const val = getNestedValue(data, candidateKey);
				if (validator(val)) return candidateKey;
			}
		}
	}
	for (const key of availableKeys) {
		const val = getNestedValue(data, key);
		if (validator(val)) return key;
	}
	return '';
}

export function findUrlMatch(data: any, availableKeys: string[], defaultKeywords: string[] = []): string {
	return findValidatedMatch(data, availableKeys, defaultKeywords, isValidUrl);
}

export function findMethodMatch(data: any, availableKeys: string[], defaultKeywords: string[] = []): string {
	return findValidatedMatch(data, availableKeys, defaultKeywords, isValidMethod);
}

export function findHeadersMatch(data: any, availableKeys: string[], defaultKeywords: string[] = []): string {
	return findValidatedMatch(data, availableKeys, defaultKeywords, isValidHeaders);
}

// --- Schema Detectors ---

/**
 * Schema: Structured HTTP Log
 * Matches logs with host + path + method fields (e.g. echo/zerolog style logs).
 * Combines host + path into a full URL.
 *
 * Example:
 * { "host": "api.example.com", "path": "/api/v1/foo", "method": "POST",
 *   "body": "{...}", "user_agent": "axios/1.6.2", ... }
 */
function tryParseStructuredHttpLog(data: any): ImportResult | null {
	const host: string = data.host;
	const path: string = data.path;
	const method: string = data.method;

	if (!host || !path || !method) return null;
	if (!isValidMethod(method)) return null;

	// Build URL: assume https if not specified
	const url = `https://${host}${path}`;

	const headers: Header[] = [];

	if (data.user_agent) {
		headers.push({ key: 'User-Agent', value: String(data.user_agent) });
	}

	// Body: may be a JSON string or object
	let body: string | undefined;
	if (data.body) {
		const rawBody = data.body;
		if (typeof rawBody === 'string') {
			try {
				const parsed = JSON.parse(rawBody);
				body = JSON.stringify(parsed, null, 2);
			} catch {
				body = rawBody;
			}
		} else if (typeof rawBody === 'object') {
			body = JSON.stringify(rawBody, null, 2);
		}
	}
	return {
		url,
		method: method.toUpperCase(),
		headers,
		body,
		metadata: { bodyType: inferBodyType(body) },
		config: { url: 'host+path', method: 'method', headers: 'user_agent', body: 'body' },
		availableKeys: extractKeyPaths(data)
	};
}

// --- Default JSON Mappings Fallback ---

export const DEFAULT_MAPPINGS = {
	url: ['url', 'path', 'endpoint', 'request.url'],
	method: ['method', 'type', 'request.method'],
	headers: ['headers', 'request.headers', 'head'],
	body: ['body', 'data', 'payload', 'request.body']
};

export function importFromJson(
	jsonInput: string,
	config: Partial<ImportConfig> = {},
	defaultMappings: any = DEFAULT_MAPPINGS
): ImportResult {
	let data: any;
	try {
		data = JSON.parse(jsonInput);
	} catch (e) {
		throw new Error('Invalid JSON input');
	}

	// Try schema-specific parsers first
	const schemas = [tryParseStructuredHttpLog];
	for (const trySchema of schemas) {
		const result = trySchema(data);
		if (result) return result;
	}

	// Fallback: generic key-path mapping
	const result: ImportResult = {
		config: { url: '', method: '', headers: '', body: '', ...config },
		availableKeys: []
	};

	try {
		const availableKeys = extractKeyPaths(data);
		result.availableKeys = availableKeys;

		if (!result.config.url) result.config.url = findUrlMatch(data, availableKeys, defaultMappings.url);
		if (!result.config.method) result.config.method = findMethodMatch(data, availableKeys, defaultMappings.method);
		if (!result.config.headers) result.config.headers = findHeadersMatch(data, availableKeys, defaultMappings.headers);
		if (!result.config.body) result.config.body = findBestMatch(availableKeys, defaultMappings.body);

		if (result.config.url) {
			const val = resolveValue(data, result.config.url);
			if (val !== undefined && val !== null) result.url = String(val);
		}

		if (result.config.method) {
			const val = resolveValue(data, result.config.method);
			if (isValidMethod(val)) result.method = String(val).toUpperCase();
		} else {
			result.method = 'GET';
		}

		if (result.config.body) {
			const val = getNestedValue(data, result.config.body);
			if (val) {
				result.body = typeof val === 'object' ? JSON.stringify(val, null, 2) : String(val);
			} else {
				const tplVal = resolveValue(data, result.config.body);
				if (tplVal) result.body = String(tplVal);
			}
		}

		if (result.config.headers) {
			const val = getNestedValue(data, result.config.headers);
			if (val && typeof val === 'object') {
				if (Array.isArray(val)) {
					result.headers = val
						.map((h: any) => ({
							key: h.key || h.name || Object.keys(h)[0] || '',
							value: h.value || Object.values(h)[0] || ''
						}))
						.filter(h => h.key);
				} else {
					result.headers = Object.entries(val).map(([k, v]) => ({ key: k, value: String(v) }));
				}
			}
		}

		// Auto-append User-Agent header if found in root
		const uaKey = findBestMatch(availableKeys, ['user_agent', 'user-agent']);
		if (uaKey) {
			const uaVal = getNestedValue(data, uaKey);
			if (typeof uaVal === 'string') {
				if (!result.headers) result.headers = [];
				const exists = result.headers.some(h => h.key.toLowerCase() === 'user-agent');
				if (!exists) result.headers.push({ key: 'User-Agent', value: uaVal });
			}
		}

		if (result.body && typeof result.body === 'string') {
			try {
				const trimmed = result.body.trim();
				if (trimmed.startsWith('{') || trimmed.startsWith('[')) {
					result.body = JSON.stringify(JSON.parse(trimmed), null, 2);
				}
			} catch {}
		}

		// Set metadata.bodyType based on whether body was detected
		result.metadata = { bodyType: inferBodyType(result.body) };

		if (!result.url) {
			throw new Error('Could not find a valid URL in the JSON mapping.');
		}
	} catch (e: any) {
		console.error('Import error', e);
		throw e;
	}

	return result;
}

// --- Parser Entry ---

function toReplay(result: ImportResult): Partial<Replay> {
	const headersObj: Record<string, string> = {};
	if (result.headers) {
		result.headers.forEach(h => (headersObj[h.key] = h.value));
	}
	const bodyType: BodyTypeHttp = result.metadata?.bodyType ?? inferBodyType(result.body);
	return {
		name: 'Imported Request',
		url: result.url || '',
		method: result.method || 'GET',
		headers: JSON.stringify(headersObj),
		payload: result.body || '',
		config: JSON.stringify({ auth: { type: 'none', config: {} }, settings: {} }),
		metadata: JSON.stringify({ params: [], bodyType })
	};
}

export const jsonParser = {
	parse: (text: string): { parsed: Partial<Replay>; importType: 'json'; rawText: string } => {
		const result = importFromJson(text);
		return {
			parsed: toReplay(result),
			importType: 'json',
			rawText: text
		};
	}
};
