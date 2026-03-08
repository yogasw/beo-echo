import { describe, it, expect } from 'vitest';
import { 
    extractKeyPaths, 
    getNestedValue, 
    findBestMatch, 
    jsonParser 
} from './jsonLogs';

describe('JSON Logs Parser Utils', () => {
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
});

describe('JSON Logs Parser', () => {
    it('should parse valid json correctly', () => {
        const json = JSON.stringify({
            request: {
                url: 'http://example.com',
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: {
                    foo: 'bar'
                }
            }
        });

        const result = jsonParser.parse(json);
        expect(result.importType).toBe('json');
        expect(result.parsed.url).toBe('http://example.com');
        expect(result.parsed.method).toBe('POST');
        expect(result.parsed.payload).toContain('"foo": "bar"');
        expect(result.parsed.headers).toContain('Content-Type');
    });

    it('should throw if no url is found', () => {
        const json = JSON.stringify({
            data: 'test'
        });
        expect(() => jsonParser.parse(json)).toThrowError('Could not find a valid URL in the JSON mapping.');
    });
});
