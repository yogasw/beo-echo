import { getLocalStorage, removeAuthLocalStorage } from '$lib/utils/localStorage';
import axios from 'axios';
import { goto } from '$app/navigation';

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
	active_proxy_id: null;
	active_proxy: null;
	endpoints: Endpoint[];
	proxy_targets: null;
	created_at: Date;
	updated_at: Date;
	url: string;
	alias: string;
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
	baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:3600/mock/api',
	headers: {
		'Content-Type': 'application/json'
	}
});

// Add request interceptor to add auth header
api.interceptors.request.use(
	(config) => {
		const username = getLocalStorage('username');
		const password = getLocalStorage('password');
		if (config.headers) {
			config.headers.Authorization = `Basic ${btoa(`${username}:${password}`)}`;
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
			removeAuthLocalStorage();
			if (!isRedirectingToLogin) {
				isRedirectingToLogin = true;
				setTimeout(() => {
					isRedirectingToLogin = false;
				}, 2000);
				if (window.location.pathname !== '/login') {
					goto('/login');
				}
			}
		}
		return Promise.reject(error);
	}
);

export const getMockStatus = async (): Promise<ConfigResponse[]> => {
	const response = await api.get('/status');
	return response.data.data;
};

export const getProjects = async (): Promise<Project[]> => {
	const response = await api.get('/projects');
	return response.data.data;
};

export const deleteProject = async (projectId: string): Promise<any> => {
	const response = await api.delete(`/projects/${projectId}`);
	return response.data;
};
export const deleteEndpoint = async (projectId: string, endpointId: string): Promise<any> => {
	const response = await api.delete(`/projects/${projectId}/endpoints/${endpointId}`);
	return response.data;
}
export const deleteResponse = async (projectId: string, endpointId: string, responseId: string): Promise<any> => {
	const response = await api.delete(`/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}`);
	return response.data;
}


export const addProject = async (name: string, alias: string): Promise<Project> => {
	const response = await api.post('/projects', {
		name,
		alias
	});
	return response.data.data;
};

export const addEndpoint = async (projectId: string, method: string, path: string): Promise<Endpoint> => {
	const response = await api.post(`/projects/${projectId}/endpoints`, {
		method,
		path,
		enabled: true,
		responseMode: 'static',
	});
	return response.data.data;
};

export const addResponse = async (projectId: string, endpointId: string, statusCode: number, body: string, headers: string): Promise<Response> => {
	const response = await api.post(`/projects/${projectId}/endpoints/${endpointId}/responses`, {
		statusCode,
		body,
		headers,
		priority: 0,
		delayMs: 0,
		stream: false,
		enabled: true,
		documentation: '',
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

export const getProjectDetail = async (uuid: string): Promise<Project> => {
	const response = await api.get(`/projects/${uuid}`);
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
	const response = await api.post('/auth', credentials);
	if (response.data.success) {
		localStorage.setItem('auth', JSON.stringify({
			username: credentials.username,
			password: credentials.password
		}));
	}
	return response.data.success;
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
