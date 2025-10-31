import apiClient from './apiClient';

export interface GenerateRequest {
	template: string;
	context?: string;
}

export interface GenerateResponse {
	content: string;
	model: string;
	token_used: number;
}

/**
 * Generate content using AI based on template
 */
export async function generateContent(request: GenerateRequest): Promise<GenerateResponse> {
	const response = await apiClient.post('/ai/generate', request);
	return response.data;
}
