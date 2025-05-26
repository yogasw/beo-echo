
export interface Replay {
	id: string;
	project_id: string;
	name: string;
	protocol: string;
	method: string;
	url: string;
	headers?: Record<string, string>;
	payload?: string;
	body?: string; // Alias for payload for compatibility
	created_at: string;
	updated_at: string;
	folder_id?: string;
	service?: string;
	method_name?: string;
	metadata?: string;
	is_mutation?: boolean;
	path?: string[];
}

export interface CreateReplayRequest {
	name: string;
	protocol: string;
	method: string;
	url: string;
	headers?: Record<string, string>;
	payload?: string;
	body?: string; // Alias for payload for compatibility
	folder_id?: string;
	service?: string;
	method_name?: string;
	metadata?: string;
	is_mutation?: boolean;
}

export interface ReplayExecutionResult {
	replay_id: string;
	status_code: number;
	response_headers?: Record<string, string>;
	response_body?: string;
	latency_ms: number;
	error?: string;
	log_id: string;
}

export interface ReplayLog {
	id: string;
	project_id: string;
	method: string;
	path: string;
	query_params?: string;
	request_headers?: string;
	request_body?: string;
	response_status: number;
	status_code: number; // Alias for response_status for compatibility
	response_body?: string;
	response_headers?: string;
	latency_ms: number;
	duration_ms: number; // Alias for latency_ms for compatibility
	response_size: number;
	bookmark?: boolean;
	logs_hash?: string;
	source: string;
	execution_mode: string;
	matched?: boolean;
	created_at: string;
	executed_at: string; // Alias for created_at for compatibility
	error_message?: string;
}

export interface ListReplaysResponse {
	replays: Replay[];
	count: number;
}

export interface ListReplayLogsResponse {
	logs: ReplayLog[];
	count: number;
}

export interface ExecuteReplayResponse {
	replay_id: string;
	status_code: number;
	response_body: string;
	response_headers: Record<string, string>;
	latency_ms: number;
	duration_ms: number; // Alias for latency_ms for compatibility
	response_size: number;
	error?: string;
	error_message?: string; // Alias for error for compatibility
	log_id: string;
	executed_at: string;
}

// HTTP Methods
export const HTTP_METHODS = [
	'GET',
	'POST',
	'PUT',
	'PATCH',
	'DELETE',
	'HEAD',
	'OPTIONS'
] as const;

export type HttpMethod = typeof HTTP_METHODS[number];

// Protocols
export const PROTOCOLS = ['http', 'https'] as const;
export type Protocol = typeof PROTOCOLS[number];
