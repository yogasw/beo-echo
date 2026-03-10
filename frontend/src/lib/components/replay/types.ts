
import type { Replay, ReplayFolder } from '$lib/types/Replay';
import type { ExecuteReplayResponse } from '$lib/types/Replay';

// Tab management with extended content storage
export interface Tab {
    id: string;
    isUnsaved: boolean;
    itemType?: 'request' | 'folder';
    folder?: ReplayFolder; // The whole folder object (UI state)

    // Replay data — null/undefined means unsaved placeholder tab
    replay?: Replay;

    // UI state
    activeSection?: string;

    // Extended content for persistence
    content?: Replay | ReplayFolder | any;

    // Per-tab execution result (null = never executed)
    executionResult?: ExecuteReplayResponse | null;
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