import { describe, it, expect } from 'vitest';
import { parseImportText } from './index';

describe('Parser Index', () => {
    it('should route to curl parser', () => {
        const result = parseImportText("curl 'http://example.com'");
        expect(result.importType).toBe('curl');
        expect(result.parsed.url).toBe('http://example.com');
    });

    it('should route to json parser', () => {
        const json = JSON.stringify({
            url: 'http://test.com',
            method: 'GET'
        });
        const result = parseImportText(json);
        expect(result.importType).toBe('json');
        expect(result.parsed.url).toBe('http://test.com');
    });

    it('should throw for unsupported formats', () => {
        expect(() => parseImportText("Random text that is not valid JSON or cURL")).toThrowError('Format not supported');
    });
});
