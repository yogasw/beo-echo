export function parseHeaders(
	headersJson?: string
): Array<{ key: string; value: string; description: string; enabled: boolean }> {
	if (!headersJson) {
		return [{ key: '', value: '', description: '', enabled: true }];
	}

	try {
		const headers = JSON.parse(headersJson);
		if (Array.isArray(headers)) {
			if (headers.length > 0) {
				return headers.map((h: any) => ({
					key: h.key || '',
					value: h.value || '',
					description: h.description || '',
					enabled: true
				}));
			}
			return [{ key: '', value: '', description: '', enabled: true }];
		}
		// Handle object format: Record<string, string> from backend
		if (typeof headers === 'object' && headers !== null) {
			const entries = Object.entries(headers);
			if (entries.length > 0) {
				return entries.map(([key, value]) => ({
					key,
					value: String(value),
					description: '',
					enabled: true
				}));
			}
		}
		return [{ key: '', value: '', description: '', enabled: true }];
	} catch (e) {
		console.error('Failed to parse headers JSON:', e);
		return [{ key: '', value: '', description: '', enabled: true }];
	}
}

export function parseConfig(configJson?: string): {
	parsedAuth?: { type: string; config: any };
	parsedSettings?: any;
} {
	if (!configJson) {
		return {
			parsedAuth: { type: 'none', config: {} },
			parsedSettings: {}
		};
	}

	try {
		const config = JSON.parse(configJson);
		return {
			parsedAuth: config.auth || { type: 'none', config: {} },
			parsedSettings: config.settings || {}
		};
	} catch (e) {
		console.error('Failed to parse config JSON:', e);
		return {
			parsedAuth: { type: 'none', config: {} },
			parsedSettings: {}
		};
	}
}

export function parseMetadata(metadataJson?: string): {
	parsedParams?: Array<{ key: string; value: string; description: string; enabled: boolean }>;
	parsedBodyType?: string;
} {
	if (!metadataJson) {
		return {
			parsedParams: [{ key: '', value: '', description: '', enabled: true }],
			parsedBodyType: 'none'
		};
	}

	try {
		const metadata = JSON.parse(metadataJson);
		return {
			parsedParams: metadata.params || [{ key: '', value: '', description: '', enabled: true }],
			parsedBodyType: metadata.bodyType || 'none'
		};
	} catch (e) {
		console.error('Failed to parse metadata JSON:', e);
		return {
			parsedParams: [{ key: '', value: '', description: '', enabled: true }],
			parsedBodyType: 'none'
		};
	}
}

export function parseUrlParams(
	url: string,
	existingParams?: Array<{ key: string; value: string; description: string; enabled: boolean }>
): Array<{ key: string; value: string; description: string; enabled: boolean }> {
	try {
		const [, queryString] = (url || '').split('?');
		const paramsList: Array<{ key: string; value: string; description: string; enabled: boolean }> = [];

		const existingKeyMap = new Map();
		if (existingParams) {
			existingParams.forEach((p) => {
				if (p.key) existingKeyMap.set(p.key, p);
			});
		}

		const usedKeys = new Set<string>();

		if (queryString) {
			const urlParams = new URLSearchParams(queryString);
			for (const [key, value] of urlParams.entries()) {
				const existing = existingKeyMap.get(key);
				paramsList.push({
					key,
					value,
					description: existing ? existing.description : '',
					enabled: true
				});
				usedKeys.add(key);
			}
		}

		// Always keep existing unused parameters (keep their description), just mark them as disabled if they're not in the URL
		if (existingParams) {
			existingParams.forEach((p) => {
				if (p.key && !usedKeys.has(p.key)) {
					paramsList.push({ ...p, enabled: false });
				}
			});
		}

		paramsList.push({ key: '', value: '', description: '', enabled: true });
		return paramsList;
	} catch (e) {
		return [{ key: '', value: '', description: '', enabled: true }];
	}
}

export function getUrlFromParams(
	baseUrl: string,
	params: Array<{ key: string; value: string; enabled: boolean }>
): string {
	try {
		const [base] = (baseUrl || '').split('?');
		const searchParams = new URLSearchParams();
		let hasParams = false;
		params.forEach((p) => {
			if (p.enabled && p.key.trim() !== '') {
				searchParams.append(p.key.trim(), p.value);
				hasParams = true;
			}
		});
		return hasParams ? `${base}?${searchParams.toString()}` : base;
	} catch (e) {
		return baseUrl || '';
	}
}
