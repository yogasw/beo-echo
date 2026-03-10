import { describe, it, expect } from 'vitest';
import { curlParser } from './curl';

describe('cURL Parser', () => {
    it('should parse basic GET request', () => {
        const result = curlParser.parse("curl 'http://example.com'");
        expect(result.importType).toBe('curl');
        expect(result.parsed.url).toBe('http://example.com');
        expect(result.parsed.method).toBe('GET');
    });

    it('should parse POST request with headers and data', () => {
        const cmd = `curl -X POST http://example.com/api -H "Content-Type: application/json" -d '{"hello": "world"}'`;
        const result = curlParser.parse(cmd);
        expect(result.parsed.method).toBe('POST');
        expect(result.parsed.url).toBe('http://example.com/api');
        expect(result.parsed.headers).toContain('Content-Type');
        expect(result.parsed.payload).toBe('{"hello": "world"}');
    });

    it('should throw for invalid command', () => {
        expect(() => curlParser.parse("not a curl command")).toThrowError('Invalid cURL command');
    });

    it('should throw if no url found', () => {
        expect(() => curlParser.parse("curl -X POST")).toThrowError('Could not find a valid command or URL in the cURL structure.');
    });
});
