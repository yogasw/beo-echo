export interface Replay {
	id: string
	name: string
	project_id: string
	folder_id: any
	protocol: string
	method: string
	url: string
	config: string
	metadata: string
	headers: string
	payload: string
	created_at: string
	updated_at: string
}

export interface CreateReplayRequest {
	name: string
	protocol: string
	method: string
	url: string
	headers?: Record<string, string>
	payload?: string
	body?: string
	folder_id?: string;
}

export interface UpdateReplayRequest {
	name?: string;
	protocol?: string;
	method?: string;
	url?: string;
	headers?: Record<string, string>;
	payload?: string;
	body?: string; // Alias for payload for compatibility
	folder_id?: string;
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
	status_text: string;
	response_body: string;
	latency_ms: number;
	size: number;
	log_id: string;
	error: string | null;
	response_headers: Record<string, string>;
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
