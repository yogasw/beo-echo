// Body type for HTTP request body
export type BodyTypeHttp = 'none' | 'raw' | 'form-data' | 'x-www-form-urlencoded';

// Typed shape of the metadata JSON field stored in Replay
export interface ReplayMetadata {
	params?: Array<{ key: string; value: string; description: string; enabled: boolean }>;
	bodyType?: BodyTypeHttp;
}

// Typed shape of the config JSON field stored in Replay
export interface ReplayConfig {
	auth?: {
		type: 'none' | 'basic' | 'bearer' | 'apiKey';
		config: Record<string, string>;
	};
	settings?: {
		timeout?: number;
		followRedirects?: boolean;
		verifySsl?: boolean;
		[key: string]: unknown;
	};
}

// Typed shape of the payload JSON field stored in Replay (for HTTP protocol)
export interface ReplayPayload {
	content: string;
	bodyType?: ReplayMetadata['bodyType'];
}

// Full Replay model — returned by GET /replays/:id
export interface Replay {
	id: string
	name: string
	doc: string
	project_id: string
	folder_id: any
	parent_id?: string
	protocol: string
	method: string
	url: string
	config: string    // JSON string
	metadata: string  // JSON string → ReplayMetadata
	headers: string   // JSON string → Record<string, string>
	payload: string

	// Response snapshot fields (only present when is_response = true)
	is_response?: boolean
	response_status?: number
	response_meta?: string
	response_body?: string
	latency_ms?: number
	created_at: string
	updated_at: string
}

// Lightweight item — returned by GET /replays (list). Heavy fields excluded.
export interface ReplayListItem {
	id: string
	name: string
	project_id: string
	folder_id: string | null | undefined
	parent_id: string | null
	is_response: boolean
	method: string
	created_at: string
	updated_at: string
}

export interface CreateReplayRequest {
	name: string
	doc?: string
	protocol: string
	method: string
	url: string
	headers?: Record<string, string>
	payload?: string
	body?: string
	folder_id?: string
	metadata?: ReplayMetadata
	config?: ReplayConfig
}

export interface UpdateReplayRequest {
	name?: string;
	doc?: string;
	protocol?: string;
	method?: string;
	url?: string;
	headers?: Record<string, string>;
	payload?: string;
	body?: string;
	folder_id?: string | null;
	update_folder_id?: boolean;
	metadata?: ReplayMetadata;
	config?: ReplayConfig;
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

export interface ReplayFolder {
	id: string;
	name: string;
	project_id: string;
	parent_id?: string;
	created_at: string;
	updated_at: string;
}

export interface ListReplaysResponse {
	replays: ReplayListItem[];
	folders: ReplayFolder[];
	replay_count: number;
	folder_count: number;
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
