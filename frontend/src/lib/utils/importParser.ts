import type { Replay } from '$lib/types/Replay';

// --- JSON Parser Utils ---

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
    
    // Convert bracket notation to dot notation
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

// --- Import Logic ---

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
  config: ImportConfig;
  availableKeys: string[];
}

export const DEFAULT_MAPPINGS = {
	url: ["url", "path", "endpoint", "request.url"],
	method: ["method", "type", "request.method"],
	headers: ["headers", "request.headers", "head"],
	body: ["body", "data", "payload", "request.body"]
};

export const METHODS = ["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"];

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
    let hasReplacement = false;
    return configStr.replace(/(\$([a-zA-Z0-9_.[\]]+))|(\$\{([a-zA-Z0-9_.[\]]+)\})/g, (match, p1, p2, p3, p4) => {
        const key = p2 || p4; 
        const val = getNestedValue(data, key);
        if (val !== undefined) {
            hasReplacement = true;
            return String(val);
        }
        return match;
    });
  }

  return undefined;
}

export function findBestMatch(availableKeys: string[], keywords: string[]): string {
  if (!keywords) return "";
  const LowerKeys = availableKeys.map(k => k.toLowerCase());
  for (const keyword of keywords) {
    if (availableKeys.includes(keyword)) return keyword;
    const index = LowerKeys.indexOf(keyword.toLowerCase());
    if (index !== -1) return availableKeys[index];
    const suffixMatch = availableKeys.find(k => k.toLowerCase().endsWith(`.${keyword.toLowerCase()}`));
    if (suffixMatch) return suffixMatch;
  }
  return "";
}

function findValidatedMatch(
  data: any, 
  availableKeys: string[], 
  defaultKeywords: string[], 
  validator: (val: any) => boolean
): string {
  if (!defaultKeywords) return "";
  for (const keyword of defaultKeywords) {
      let val: any = undefined;
      let candidateKey = "";

      if (keyword.includes('$')) {
         val = resolveValue(data, keyword);
         if (val !== undefined && validator(val)) {
             return keyword; 
         }
      } else {
         candidateKey = findBestMatch(availableKeys, [keyword]);
         if (candidateKey) {
             val = getNestedValue(data, candidateKey);
             if (validator(val)) {
                 return candidateKey;
             }
         }
      }
  }

  for (const key of availableKeys) {
    const val = getNestedValue(data, key);
    if (validator(val)) return key;
  }
  
  return "";
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

export function importFromJson(
  jsonInput: string, 
  config: Partial<ImportConfig> = {}, 
  defaultMappings: any = DEFAULT_MAPPINGS
): ImportResult {
  const result: ImportResult = { 
      config: { url: '', method: '', headers: '', body: '', ...config },
      availableKeys: []
  };
  
  try {
    const data = JSON.parse(jsonInput);
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
      if (isValidMethod(val)) {
        result.method = String(val).toUpperCase();
      }
    } else {
		result.method = "GET"; // default fallback
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
           result.headers = val.map((h: any) => ({
             key: h.key || h.name || Object.keys(h)[0] || "",
             value: h.value || Object.values(h)[0] || ""
           })).filter(h => h.key);
        } else {
          result.headers = Object.entries(val).map(([k, v]) => ({ key: k, value: String(v) }));
        }
      }
    }

    const uaKey = findBestMatch(availableKeys, ['user_agent', 'user-agent', 'user-agent-data']);
    if (uaKey) {
       const uaVal = getNestedValue(data, uaKey);
       if (typeof uaVal === 'string') {
          if (!result.headers) result.headers = [];
          const exists = result.headers.some(h => h.key.toLowerCase() === 'user-agent');
          if (!exists) {
             result.headers.push({ key: 'User-Agent', value: uaVal });
          }
       }
    }

    if (result.body && typeof result.body === 'string') {
        try {
            const trimmed = result.body.trim();
            if (trimmed.startsWith('{') || trimmed.startsWith('[')) {
                 const parsed = JSON.parse(trimmed);
                 result.body = JSON.stringify(parsed, null, 2);
            }
        } catch (e) {}
    }

	if (!result.url) {
		throw new Error("Could not find a valid URL in the JSON mapping.");
	}

  } catch (e: any) {
    console.error("Import error", e);
    throw e;
  }

  return result;
}

// Basic cURL parser to extract URL, Method, Headers, and Body
export function parseCurl(curlText: string): Partial<Replay> {
    const lines = curlText.match(/(?:[^\s']+|'[^']*')+/g) || [];
    if (lines.length === 0 || !lines[0]?.toLowerCase().includes("curl")) {
        throw new Error("Invalid cURL command");
    }

    let url = "";
    let method = "GET";
    let headers: Array<{key: string, value: string}> = [];
    let body = "";

    let skipNext = false;

    for (let i = 1; i < lines.length; i++) {
        let token = lines[i];
		
		if (token === '\\') continue;
		
		if (skipNext) {
			skipNext = false;
			continue;
		}

		// remove surrounding quotes if exist
		if ((token.startsWith("'") && token.endsWith("'")) || (token.startsWith('"') && token.endsWith('"'))) {
			token = token.substring(1, token.length - 1);
		}

        if (token === "-X" || token === "--request") {
            method = lines[i+1]?.replace(/['"]/g, '').toUpperCase() || method;
			skipNext = true;
        } else if (token === "-H" || token === "--header") {
            const headerStr = lines[i+1]?.replace(/['"]/g, '') || "";
            const separatorIdx = headerStr.indexOf(':');
            if (separatorIdx > -1) {
                headers.push({
                    key: headerStr.substring(0, separatorIdx).trim(),
                    value: headerStr.substring(separatorIdx + 1).trim()
                });
            }
			skipNext = true;
        } else if (token === "-d" || token === "--data" || token === "--data-raw" || token === "--data-binary" || token === "--data-urlencode") {
            body = lines[i+1] || "";
			if ((body.startsWith("'") && body.endsWith("'")) || (body.startsWith('"') && body.endsWith('"'))) {
				body = body.substring(1, body.length - 1);
			}
            if (method === "GET") method = "POST";
			skipNext = true;
        } else if (!token.startsWith('-') && !url) {
            url = token;
        } else if (token.startsWith('-')) {
			// Skip next token if it's a known flag that takes an argument but we don't handle it
			if (["-u", "--user", "-A", "--user-agent", "-e", "--referer", "-m", "--max-time", "-w", "--write-out", "-x", "--proxy", "-b", "--cookie", "-c", "--cookie-jar", "-F", "--form"].includes(token)) {
				skipNext = true;
			}
		}
    }

	if (!url) {
		throw new Error("Could not find a valid command or URL in the cURL structure.");
	}

	const headersObj: any = {};
	headers.forEach(h => {
		headersObj[h.key] = h.value;
	});

    return {
        name: "Imported Request",
        url,
        method,
        headers: JSON.stringify(headersObj),
        payload: body,
		config: JSON.stringify({ auth: { type: 'none', config: {} }, settings: {} })
    };
}

export function parseImportText(text: string): { parsed: Partial<Replay>; importType: 'curl' | 'json', rawText: string } {
	const trimmed = text.trim();
	if (trimmed.toLowerCase().startsWith('curl')) {
		return {
			parsed: parseCurl(trimmed),
			importType: 'curl',
			rawText: trimmed
		};
	} else if (trimmed.startsWith('{') || trimmed.startsWith('[')) {
		const result = importFromJson(trimmed);
		
		const headersObj: any = {};
		if (result.headers) {
			result.headers.forEach((h: any) => {
				headersObj[h.key] = h.value;
			});
		}

		return {
			parsed: {
				name: "Imported Request",
				url: result.url || "",
				method: result.method || "GET",
				headers: JSON.stringify(headersObj),
				payload: result.body || "",
				config: JSON.stringify({ auth: { type: 'none', config: {} }, settings: {} })
			},
			importType: 'json',
			rawText: trimmed
		};
	} else {
		throw new Error("Format not supported. Please paste a valid JSON or cURL command.");
	}
}
