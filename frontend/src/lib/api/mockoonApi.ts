import axios from 'axios';
import { goto } from '$app/navigation';
import type { User } from '$lib/types/User';
import { auth } from '$lib/stores/auth';
import { BASE_URL_API } from '$lib/utils/authUtils';

interface AuthCredentials {
	username: string;
	password: string;
}

export interface ConfigResponse {
	uuid: string;
	name: string;
	configFile: string;
	port: number;
	url: string;
	size: string;
	modified: string;
	inUse: boolean;
}

export type Project = {
	id: string;
	name: string;
	mode: string;
	status?: string; // 'running', 'stopped', or 'error'
	active_proxy_id: null;
	active_proxy: null;
	endpoints: Endpoint[];
	proxy_targets: null;
	created_at: Date;
	updated_at: Date;
	url: string;
	alias: string;
}

export type RequestLog = {
	id: string;
	project_id: string;
	method: string;
	path: string;
	query_params: string;
	request_headers: string;
	request_body: string;
	response_status: number;
	response_body: string;
	response_headers: string;
	latency_ms: number;
	execution_mode: string;
	matched: boolean;
	created_at: Date;
}

export type Endpoint = {
	id: string;
	project_id: string;
	method: string;
	path: string;
	enabled: boolean;
	response_mode: string;
	responses: Response[];
	created_at: Date;
	updated_at: Date;
	documentation: string;
}

export type Response = {
	id: string;
	endpoint_id: string;
	status_code: number;
	body: string;
	headers: string;
	priority: number;
	delay_ms: number;
	stream: boolean;
	enabled: boolean;
	documentation: string;
	rules: null;
	created_at: Date;
	updated_at: Date;
}

// Create axios instance with default config
const api = axios.create({
	baseURL: BASE_URL_API,
	headers: {
		'Content-Type': 'application/json'
	}
});

// Add request interceptor to add auth header with JWT token
api.interceptors.request.use(
	(config) => {
		// Get JWT token from auth store
		const token = auth.getToken();
		
		if (config.headers && token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

let isRedirectingToLogin = false;
// Add response interceptor for handling auth errors
api.interceptors.response.use(
	response => response,
	error => {
		console.log('route', error.response?.config.url);
		if (error.response?.status === 401) {
			console.error('Authentication failed');
			// Use auth store logout to properly clear token and state
			auth.logout();
			// Prevent multiple redirects
			if (!isRedirectingToLogin) {
				isRedirectingToLogin = true;
				setTimeout(() => {
					isRedirectingToLogin = false;
				}, 2000);
			}
		}
		return Promise.reject(error);
	}
);

export const getMockStatus = async (): Promise<ConfigResponse[]> => {
	const response = await api.get('/status');
	return response.data.data;
};

export const getWorkspaces = async (): Promise<any[]> => {
	const response = await api.get('/workspaces');
	return response.data.data;
};

export const getProjects = async (workspaceId?: string): Promise<Project[]> => {
	// If no workspaceId provided, this will fail with the new API structure
	if (!workspaceId) {
		console.error('No workspace ID provided for getProjects');
		return [];
	}
	const response = await api.get(`/workspaces/${workspaceId}/projects`);
	return response.data.data;
};

export const deleteProject = async (workspaceId: string, projectId: string): Promise<any> => {
	const response = await api.delete(`/workspaces/${workspaceId}/projects/${projectId}`);
	return response.data;
};
export const deleteEndpoint = async (workspaceId: string, projectId: string, endpointId: string): Promise<any> => {
	const response = await api.delete(`/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}`);
	return response.data;
}
export const deleteResponse = async (workspaceId: string, projectId: string, endpointId: string, responseId: string): Promise<any> => {
	const response = await api.delete(`/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}`);
	return response.data;
}

export const updateProjectStatus = async (workspaceId: string, projectId: string, status: string): Promise<Project> => {
	const response = await api.put(`/workspaces/${workspaceId}/projects/${projectId}`, {
		status: status
	});
	return response.data.data;
};

export const addProject = async (workspaceId: string, name: string, alias: string): Promise<Project> => {
	const response = await api.post(`/workspaces/${workspaceId}/projects`, {
		name,
		alias
	});
	return response.data.data;
};

export const addEndpoint = async (workspaceId: string, projectId: string, method: string, path: string): Promise<Endpoint> => {
	const response = await api.post(`/workspaces/${workspaceId}/projects/${projectId}/endpoints`, {
		method,
		path,
		enabled: true,
		response_mode: 'random',
	});
	return response.data.data;
};

export const addResponse = async (workspaceId: string, projectId: string, endpointId: string, statusCode: number, body: string, headers: string, documentation: string): Promise<Response> => {
	const response = await api.post(`/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses`, {
		statusCode,
		body,
		headers,
		priority: 0,
		delayMs: 0,
		stream: false,
		enabled: true,
		documentation,
	});
	return response.data.data;
};

export const uploadConfig = async (formData: FormData): Promise<any> => {
	const response = await api.post('/upload', formData, {
		headers: {
			'Content-Type': 'multipart/form-data'
		}
	});
	return response.data.data;
};

export const downloadConfig = async (filename: string): Promise<any> => {
	return await api.get(`/configs/${filename}/download`);
};

export const getProjectDetail = async (workspaceId: string, projectId: string): Promise<Project> => {
	const response = await api.get(`/workspaces/${workspaceId}/projects/${projectId}`);
	return response.data.data;
}

export const startMockServer = async (port: number, configFile: string, uuid: string): Promise<any> => {
	const response = await api.post('/start', {
		uuid,
		port,
		configFile
	});
	return response.data;
};

export const stopMockServer = async (port: number): Promise<any> => {
	const response = await api.post('/stop', {
		port
	});
	return response.data;
};

export const deleteConfig = async (filename: string): Promise<any> => {
	const response = await api.delete(`/configs/${filename}`);
	return response.data;
};

export const login = async (credentials: AuthCredentials): Promise<boolean> => {
	// Note: This function is now deprecated - use auth.login from auth store instead
	// Adapting old function to use the new auth system
	try {
		await auth.login(credentials.username, credentials.password);
		return true;
	} catch (error) {
		console.error('Login failed:', error);
		return false;
	}
};

export const getConfigDetails = async (uuid: string): Promise<ConfigResponse> => {
	const response = await api.get(`/configs/${uuid}`);
	return response.data.data;
};


export const saveGitConfig = async (config: {
	gitName: string;
	gitEmail: string;
	gitBranch: string;
	sshKey: string;
	gitUrl: string;
}): Promise<{ success: boolean; message: string }> => {
	const response = await api.post('/git/save-config', config);
	return response.data;
};

export const saveAndTestSyncGit = async (config: {
	gitName: string;
	gitEmail: string;
	gitBranch: string;
	sshKey: string;
	gitUrl: string;
}): Promise<{ success: boolean; message: string }> => {
	const response = await api.post('/git/save-and-test-sync', config);
	return response.data;
};

export const getGitConfig = async (): Promise<{
	success: boolean;
	data?: {
		gitName: string;
		gitEmail: string;
		gitUrl: string;
		gitBranch: string;
		sshKey: string;
	};
	message?: string;
}> => {
	const response = await api.get('/git/config');
	return response.data;
};

export const updateEndpoint = async (projectId: string, endpointId: string, data: {
	method?: string;
	path?: string;
	enabled?: boolean;
	responseMode?: string;
	documentation?: string;
}): Promise<Endpoint> => {
	const response = await api.put(`/projects/${projectId}/endpoints/${endpointId}`, data);
	return response.data.data;
};

export const updateResponse = async (projectId: string, endpointId: string, responseId: string, data: {
	statusCode?: number;
	body?: string;
	headers?: string;
	priority?: number;
	delayMS?: number;
	stream?: boolean;
	enabled?: boolean;
}): Promise<Response> => {
	const response = await api.put(`/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}`, data);
	return response.data.data;
};

// Get logs with pagination
export const getLogs = async (page: number = 1, pageSize: number = 100, workspaceId: string, projectId: string): Promise<{ logs: RequestLog[], total: number }> => {
	const params: Record<string, string> = {
		page: page.toString(),
		pageSize: pageSize.toString()
	};

	const response = await api.get(`/workspaces/${workspaceId}/projects/${projectId}/logs`, { params });
	return {
		logs: response.data.logs,
		total: response.data.total
	};
};

// Create an EventSource for real-time log streaming
export const createLogStream = (workspaceId: string, projectId: string, limit: number = 100): EventSource => {
	let baseURL = `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:3600/mock/api'}`;
	let url = `${baseURL}/workspaces/${workspaceId}/projects/${projectId}/logs/stream?limit=${limit}`;

	// Add authentication using JWT token
	const token = auth.getToken();
	if (token) {
		url += `&auth=${token}`;
	}

	console.log('Creating SSE connection to:', url);
	
	const eventSource = new EventSource(url);
	
	// Add global error handling
	eventSource.onerror = (err) => {
		console.error('EventSource global error:', err);
	};
	
	return eventSource;
};


// Function fetchUserProfile has been moved to lib/utils/authUtils.ts
