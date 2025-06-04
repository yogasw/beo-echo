
// Tab management with extended content storage
export interface Tab {
    id: string;
    name: string;
    method: string;
    url: string;
    isUnsaved: boolean;
    
    // Extended content for persistence
    content?: TabContent;
}

// Complete tab content that needs to be persisted
export interface TabContent {
    // Basic request data
    method: string;
    url: string;
    
    // Request components
    params: Param[];
    headers: Header[];
    body: {
        type: 'none' | 'form-data' | 'x-www-form-urlencoded' | 'raw' | 'binary';
        content: string;
        formData?: Array<{ key: string; value: string; type: 'text' | 'file'; enabled: boolean; }>;
        urlEncoded?: Array<{ key: string; value: string; enabled: boolean; }>;
    };
    
    // Authentication
    auth: AuthConfig;
    
    // Scripts
    scripts: ScriptConfig;
    
    // Settings
    settings: SettingsConfig;
    
    // UI state
    activeSection: string;
}

// Data structures for tab components
export interface Param {
    key: string;
    value: string;
    description: string;
    enabled: boolean;
}

export interface Header {
    key: string;
    value: string;
    description: string;
    enabled: boolean;
}

export interface AuthConfig {
    type: string;
    config: Record<string, any>;
}

export interface ScriptConfig {
    preRequestScript: string;
    testScript: string;
}

export interface SettingsConfig {
    timeout: number;
    followRedirects: boolean;
    maxRedirects: number;
    verifySsl: boolean;
    ignoreSslErrors: boolean;
    encoding: string;
    sendCookies: boolean;
    storeCookies: boolean;
    keepAlive: boolean;
    userAgent: string;
    retryOnFailure: boolean;
    retryCount: number;
    retryDelay: number;
}