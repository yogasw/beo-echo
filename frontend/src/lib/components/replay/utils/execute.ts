export function buildExecutePayload(content: any, source?: string) {
	// Prepare request data
	const requestData: {
		method: string;
		url: string;
		headers: Record<string, string>;
		query: Record<string, string>;
		payload: string;
		source?: string;
	} = {
		method: content?.method || 'GET',
		url: content?.url || '',
		headers: {},
		query: {},
		payload: content?.body?.content || ''
	};

	// Process headers
	if (content?.headers) {
		content.headers.forEach((header: any) => {
			if (header.enabled && header.key && header.value) {
				requestData.headers[header.key] = header.value;
			}
		});
	}

	// Process query parameters
	if (content?.params) {
		content.params.forEach((param: any) => {
			if (param.enabled && param.key && param.value) {
				requestData.query[param.key] = param.value;
			}
		});
	}

	// Process auth
	if (content?.auth?.type !== 'none') {
		const authConfig = content?.auth?.config;

		switch (content?.auth?.type) {
			case 'basic':
				// Add Basic Auth header
				if (authConfig?.username) {
					const credentials = btoa(`${authConfig?.username}:${authConfig?.password || ''}`);
					requestData.headers['Authorization'] = `Basic ${credentials}`;
				}
				break;
			case 'bearer':
				// Add Bearer token header
				if (authConfig?.token) {
					requestData.headers['Authorization'] = `Bearer ${authConfig?.token}`;
				}
				break;
			case 'apiKey':
				// Add API key as header or query param
				if (authConfig?.key && authConfig?.value) {
					if (authConfig?.in === 'header') {
						requestData.headers[authConfig.key] = authConfig.value;
					} else if (authConfig?.in === 'query') {
						requestData.query[authConfig.key] = authConfig.value;
					}
				}
				break;
			// Add other auth types as needed
		}
	}

	// Attach source if provided
	if (source) {
		requestData.source = source;
	}

	return requestData;
}
