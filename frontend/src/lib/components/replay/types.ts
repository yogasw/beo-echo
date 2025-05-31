
// Tab management
export interface Tab {
    id: string;
    name: string;
    method: string;
    url: string;
    isUnsaved: boolean;
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