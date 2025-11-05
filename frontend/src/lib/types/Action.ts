// Action types and interfaces for the Actions feature

export type ActionType = 'replace_text' | 'add_header' | 'webhook' | 'delay';

export type ExecutionPoint = 'before_request' | 'after_request';

// Action Type Information
export interface ActionTypeInfo {
	id: ActionType;
	name: string;
	description: string;
	icon: string;
	category: string;
	fields: string[];
}

export interface ActionFilter {
	id?: string;
	action_id?: string;
	type: 'method' | 'path' | 'header' | 'query' | 'status_code';
	key?: string;
	operator: 'equals' | 'contains' | 'regex' | 'starts_with' | 'ends_with';
	value: string;
	created_at?: string;
	updated_at?: string;
}

export interface Action {
	id: string;
	project_id: string;
	name: string;
	type: ActionType;
	execution_point: ExecutionPoint;
	enabled: boolean;
	priority: number;
	config: string; // JSON string
	filters?: ActionFilter[];
	created_at: string;
	updated_at: string;
}

// Replace Text Action Config
export interface ReplaceTextConfig {
	target: 'request_body' | 'response_body' | 'request_header' | 'response_header';
	pattern: string;
	replacement: string;
	use_regex: boolean;
	header_key?: string; // Required for header targets
}

// API Request/Response Types
export interface CreateActionRequest {
	name: string;
	type: ActionType;
	execution_point?: ExecutionPoint;
	enabled?: boolean;
	config: string;
	filters?: Omit<ActionFilter, 'id' | 'action_id' | 'created_at' | 'updated_at'>[];
}

export interface UpdateActionRequest {
	name?: string;
	execution_point?: ExecutionPoint;
	enabled?: boolean;
	priority?: number;
	config?: string;
	filters?: Omit<ActionFilter, 'id' | 'action_id' | 'created_at' | 'updated_at'>[];
}

export interface ListActionsResponse {
	success: boolean;
	data: Action[];
}

export interface ActionResponse {
	success: boolean;
	data: Action;
	message?: string;
}

export interface ActionMessageResponse {
	success: boolean;
	message: string;
}
