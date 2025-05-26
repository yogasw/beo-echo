
// Tab management
interface Tab {
    id: string;
    name: string;
    method: string;
    url: string;
    isUnsaved: boolean;
}

// Data structures for tab components
interface Param {
    key: string;
    value: string;
    description: string;
    enabled: boolean;
}

interface Header {
    key: string;
    value: string;
    description: string;
    enabled: boolean;
}

interface AuthConfig {
    type: string;
    config: Record<string, any>;
}

interface ScriptConfig {
    preRequestScript: string;
    testScript: string;
}

interface SettingsConfig {
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