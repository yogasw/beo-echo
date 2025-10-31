import apiClient from './apiClient';

export interface GenerateRequest {
	message: string;
	context?: string;
	content_type?: string;
}

export interface GenerateResponse {
	content: string;
	model: string;
	token_used: number;
	can_apply: boolean;
	data: string;
}

/**
 * Generate content using AI based on template
 */
export async function generateContent(request: GenerateRequest): Promise<GenerateResponse> {
	const response = await apiClient.post('/ai/generate', request);
	return response.data;
}
