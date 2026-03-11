export function buildExecutePayload(content: any, source?: string) {
	// Prepare request data
	const requestData: {
		method: string;
		url: string;
		headers: Record<string, string>;
		query: Record<string, string>;
		payload: string;
		source?: string;
		metadata: Record<string, string>;
	} = {
		method: content?.method || 'GET',
		url: content?.url || '',
		headers: {},
		query: {},
		payload: content?.body?.content || '',
		metadata: {
			bodyType: content?.body?.type || 'none'
		}
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

export async function executeRequest(
	workspaceId: string,
	projectId: string,
	replayData: any
): Promise<any> { // Or import ExecuteReplayResponse if preferred, avoiding circular dep if needed
	const payload = {
		protocol: 'http', // Default to http
		method: replayData.method || 'GET',
		url: replayData.url || '',
		headers: replayData.headers || {},
		payload: replayData.payload || '',
		query: replayData.query || {},
		metadata: replayData.metadata || {}
	};

	let result: any; // ExecuteReplayResponse
	
	if (replayData.source === 'browser') {
		// Execute directly from frontend using axios
		const startTime = performance.now();
		
		// Build absolute URL including query parameters
		let requestUrl = payload.url;
		if (Object.keys(payload.query).length > 0) {
			const urlObj = requestUrl.startsWith('http') ? new URL(requestUrl) : new URL(`http://${requestUrl}`);
			Object.entries(payload.query).forEach(([key, value]) => {
				urlObj.searchParams.append(key, value as string);
			});
			requestUrl = urlObj.toString();
		}

		// Set default content type if missing based on metadata
		if (!payload.headers['Content-Type'] && !payload.headers['content-type']) {
			if (payload.metadata.bodyType === 'x-www-form-urlencoded') {
				payload.headers['Content-Type'] = 'application/x-www-form-urlencoded';
			} else if (payload.metadata.bodyType === 'raw') {
				payload.headers['Content-Type'] = 'application/json';
			}
		}

		try {
			// dynamically import axios inside function to avoid heavy bundle if not needed immediately
			const axios = (await import('axios')).default;
			const response = await axios({
				method: payload.method,
				url: requestUrl,
				headers: payload.headers,
				data: payload.payload,
				validateStatus: () => true // Resolve on any status code
			});
			
			const endTime = performance.now();
			let responseDataString = typeof response.data === 'string' ? response.data : JSON.stringify(response.data, null, 2);
			if (responseDataString === undefined) responseDataString = '';
			
			// Format response to match ExecuteReplayResponse
			result = {
				replay_id: 'local',
				status_code: response.status,
				status_text: response.statusText,
				response_body: responseDataString,
				latency_ms: Math.round(endTime - startTime),
				size: new Blob([responseDataString]).size,
				log_id: 'local',
				error: null,
				response_headers: response.headers as Record<string, string>
			};
		} catch (error: any) {
			const endTime = performance.now();
			result = {
				replay_id: 'local',
				status_code: 0,
				status_text: 'Error',
				response_body: '',
				latency_ms: Math.round(endTime - startTime),
				size: 0,
				log_id: 'local',
				error: error.message || 'Network Error',
				response_headers: {}
			};
		}
	} else {
		// Calculate imports statically if needed, but replayApi comes from stores.
		// For simplicity, we import replayApi directly.
		const { replayApi } = await import('$lib/api/replayApi');
		
		// Execute via Beo Echo backend
		result = await replayApi.executeReplayRequest(
			workspaceId, 
			projectId, 
			payload
		);
	}

	return result;
}

export function replaceUrlHostToLocalhost(originalUrl: string): string {
	try {
		if (!originalUrl) return originalUrl;
		const urlWithProtocol = originalUrl.startsWith('http') ? originalUrl : `http://${originalUrl}`;
		const urlObj = new URL(urlWithProtocol);
		urlObj.hostname = 'localhost';
		return urlObj.toString();
	} catch (e) {
		console.warn('Could not parse URL to replace host with localhost:', e);
		return originalUrl;
	}
}
