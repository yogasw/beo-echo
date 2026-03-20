import { describe, it, expect } from 'vitest';
import {
	extractKeyPaths,
	getNestedValue,
	findBestMatch,
	jsonParser,
	importFromJson
} from './jsonLogs';

// --- Utilities ---

describe('extractKeyPaths', () => {
	it('should extract paths from nested objects', () => {
		const data = { a: { b: 1 }, c: 2 };
		expect(extractKeyPaths(data)).toEqual(['a', 'a.b', 'c']);
	});

	it('should extract paths from arrays', () => {
		const data = { a: [1, { b: 2 }] };
		expect(extractKeyPaths(data)).toEqual(['a', 'a[0]', 'a[1]', 'a[1].b']);
	});
});

describe('getNestedValue', () => {
	it('should retrieve nested values', () => {
		const data = { a: { b: 1 }, c: [10, { d: 20 }] };
		expect(getNestedValue(data, 'a.b')).toBe(1);
		expect(getNestedValue(data, 'c[1].d')).toBe(20);
	});
});

describe('findBestMatch', () => {
	it('should return keyword if it exists exactly', () => {
		const keys = ['user', 'user.name', 'name'];
		expect(findBestMatch(keys, ['name'])).toBe('name');
	});

	it('should return a suffix match', () => {
		const keys = ['request.url', 'id'];
		expect(findBestMatch(keys, ['url'])).toBe('request.url');
	});
});

// --- Schema: Structured HTTP Log ---

describe('Structured HTTP Log schema (host + path + method)', () => {
	const log = {
		level: 'error',
		request_id: '8bb0599bb',
		remote_ip: '194.56.225.16',
		host: 'api.example.com',
		user_agent: 'axios/1.6.2',
		method: 'POST',
		path: '/api',
		body: JSON.stringify({
			status: 'success',
		}),
		status_code: 401,
		latency: 0.06,
		time: '2026-03-05T22:34:32Z',
		message: 'http request'
	};

	it('should build URL from host + path', () => {
		const result = importFromJson(JSON.stringify(log));
		expect(result.url).toBe('https://api.example.com/api');
	});

	it('should extract method', () => {
		const result = importFromJson(JSON.stringify(log));
		expect(result.method).toBe('POST');
	});

	it('should parse body (JSON string) into readable JSON', () => {
		const result = importFromJson(JSON.stringify(log));
		expect(result.body).toContain('"status": "success"');
	});

	it('should include User-Agent header from user_agent field', () => {
		const result = importFromJson(JSON.stringify(log));
		const ua = result.headers?.find(h => h.key === 'User-Agent');
		expect(ua?.value).toBe('axios/1.6.2');
	});

	it('should format via jsonParser.parse', () => {
		const result = jsonParser.parse(JSON.stringify(log));
		expect(result.importType).toBe('json');
		expect(result.parsed.url).toBe('https://api.example.com/api');
		expect(result.parsed.method).toBe('POST');
		const headers = JSON.parse(result.parsed.headers || '{}');
		expect(headers['User-Agent']).toBe('axios/1.6.2');
	});

	it('should work with GET method (no body)', () => {
		const getLog = { ...log, method: 'GET', body: undefined };
		const result = importFromJson(JSON.stringify(getLog));
		expect(result.method).toBe('GET');
		expect(result.body).toBeUndefined();
	});
});

// --- Schema: Generic JSON Fallback ---

describe('Generic JSON fallback', () => {
	it('should parse standard JSON with explicit url field', () => {
		const json = JSON.stringify({
			request: {
				url: 'http://example.com',
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: { foo: 'bar' }
			}
		});
		const result = jsonParser.parse(json);
		expect(result.importType).toBe('json');
		expect(result.parsed.url).toBe('http://example.com');
		expect(result.parsed.method).toBe('POST');
		expect(result.parsed.payload).toContain('"foo": "bar"');
		expect(result.parsed.headers).toContain('Content-Type');
	});

	it('should throw if no url found', () => {
		const json = JSON.stringify({ data: 'test' });
		expect(() => jsonParser.parse(json)).toThrowError('Could not find a valid URL in the JSON mapping.');
	});

	it('should parse headers correctly when they are arrays of strings', () => {
	const payload = {
		"level": "info",
		"request_id": "9f4c1e8a7b2d4c6e8f1a0b3d5c7e9f10",
		"method": "POST",
		"url": "https://api.example-service.com",
		"body": "{\"worker_id\":45291,\"limit\":3,\"task_id\":908172}",
		"headers": {
			"Accept": ["application/json"],
			"Content-Type": ["application/json"],
			"X-Client-Id": ["client-7f9a3c"],
			"X-Auth-Token": ["******"]
		},
		"status_code": 200,
		"response_body": "{\"result\":{\"status\":\"ok\"}}",
		"latency": "187.341221ms",
		"time": "2026-03-13T13:02:10Z",
		"message": "outbound request"
	};

	const result = jsonParser.parse(JSON.stringify(payload));
	
	expect(result.importType).toBe('json');
	expect(result.parsed.url).toBe('https://api.example-service.com');
	expect(result.parsed.method).toBe('POST');
	
	const headers = JSON.parse(result.parsed.headers || '{}');
	expect(headers['Accept']).toBe('application/json');
	expect(headers['Content-Type']).toBe('application/json');
	expect(headers['X-Client-Id']).toBe('client-7f9a3c');
	expect(headers['X-Auth-Token']).toBe('******');
});
});
