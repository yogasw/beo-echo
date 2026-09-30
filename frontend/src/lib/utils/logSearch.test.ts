import { describe, expect, it } from 'vitest';
import type { RequestLog } from '$lib/api/BeoApi';
import { highlightSegments, highlightTermsFor, matchLog, parseSearch } from './logSearch';

function makeLog(overrides: Partial<RequestLog> = {}): RequestLog {
	return {
		id: '1',
		project_id: 'p1',
		method: 'GET',
		path: '/organization/abc/productVariant/stock',
		query_params: 'brandId=wardah&limit=1&page=1&q=WRMKCU33NF',
		request_headers: JSON.stringify({ Accept: 'application/json', 'Cf-Ray': 'a4327777b905fdcc-SIN' }),
		request_body: '',
		response_status: 200,
		response_body: JSON.stringify({ data: [{ sku: 'SKU-001', stock: 5 }] }),
		response_headers: JSON.stringify({ 'Content-Type': 'application/json' }),
		latency_ms: 12,
		execution_mode: 'forwarder',
		matched: false,
		created_at: new Date(),
		...overrides
	};
}

describe('parseSearch', () => {
	it('keeps unknown prefixes as plain text', () => {
		const { terms } = parseSearch('brandId=wardah');
		expect(terms).toHaveLength(1);
		expect(terms[0].values).toEqual(['brandid=wardah']);
		expect(terms[0].fields).not.toContain('status');
	});

	it('supports field:value, field=value, quotes, OR and negation', () => {
		const { terms } = parseSearch('body:abc header="x y" a|b -status:500');
		expect(terms[0]).toMatchObject({ values: ['abc'], fields: ['request_body', 'response_body'] });
		expect(terms[1]).toMatchObject({ values: ['x y'], fields: ['request_headers', 'response_headers'] });
		expect(terms[2].values).toEqual(['a', 'b']);
		expect(terms[3]).toMatchObject({ values: ['500'], fields: ['status'], negate: true });
	});
});

describe('matchLog', () => {
	it('finds a value that only lives in the query string', () => {
		const result = matchLog(makeLog(), parseSearch('WRMKCU33NF'));
		expect(result.matched).toBe(true);
		expect(result.fields).toEqual(['query']);
	});

	it('searches headers as "key: value"', () => {
		expect(matchLog(makeLog(), parseSearch('cf-ray')).fields).toEqual(['request_headers']);
		expect(matchLog(makeLog(), parseSearch('res.header:"content-type: application/json"')).matched).toBe(true);
		expect(matchLog(makeLog(), parseSearch('req.header:content-type')).matched).toBe(false);
	});

	it('matches JSON fragments in the body', () => {
		expect(matchLog(makeLog(), parseSearch('res.body:"stock":5')).matched).toBe(true);
		expect(matchLog(makeLog(), parseSearch('body:sku-001')).fields).toEqual(['response_body']);
		expect(matchLog(makeLog(), parseSearch('req.body:sku-001')).matched).toBe(false);
	});

	it('requires every term and honours negation and status', () => {
		expect(matchLog(makeLog(), parseSearch('wardah missing')).matched).toBe(false);
		expect(matchLog(makeLog(), parseSearch('wardah -forbidden')).matched).toBe(true);
		expect(matchLog(makeLog(), parseSearch('wardah -status:200')).matched).toBe(false);
		expect(matchLog(makeLog(), parseSearch('status:200 method:get')).matched).toBe(true);
		// a bare number does not match on status alone
		expect(matchLog(makeLog(), parseSearch('200')).matched).toBe(false);
	});
});

describe('highlight helpers', () => {
	it('returns terms only for the fields they target', () => {
		const search = parseSearch('abc body:def -ghi');
		expect(highlightTermsFor(search, 'query')).toEqual(['abc']);
		expect(highlightTermsFor(search, 'response_body')).toEqual(['abc', 'def']);
	});

	it('splits text into case-insensitive, merged segments', () => {
		expect(highlightSegments('q=WRMKCU33NF&x', ['wrmk', 'kcu33'])).toEqual([
			{ text: 'q=', match: false },
			{ text: 'WRMKCU33', match: true },
			{ text: 'NF&x', match: false }
		]);
	});
});
