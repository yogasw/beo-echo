import { writable } from 'svelte/store';
import type { RequestLog } from '$lib/api/BeoApi';

// Initialize logs store with empty array
export const logs = writable<RequestLog[]>([]);

// Store for managing live connection status and related information
export const logsConnectionStatus = writable({
    isConnected: false,
    reconnectAttempts: 0,
    projectId: null as string | null,
    autoScroll: true,
    isLoading: false,
    error: null as string | null,
    total: 0
});

// Constants for reconnection logic
export const MAX_RECONNECT_ATTEMPTS = 5;
export const RECONNECT_DELAY_MS = 3000;

// Function to clear logs (useful when switching projects)
export function clearLogs() {
    logs.set([]);
    logsConnectionStatus.update(state => ({
        ...state,
        total: 0
    }));
}

// Function to add a new log to the store
export function addLog(newLog: RequestLog) {
    console.log('Adding log:', newLog);
    logs.update(currentLogs => {
        // Check if log already exists to prevent duplicates
        if (!currentLogs.some(log => log.id === newLog.id)) {
            // Add to beginning of array (newest first) and limit to 1000 logs to prevent browser slowdown
            return [newLog, ...currentLogs].slice(0, 1000);
        }
        return currentLogs;
    });
}

// Update connection status
export function updateConnectionStatus(status: boolean) {
    logsConnectionStatus.update(state => ({
        ...state,
        isConnected: status,
        reconnectAttempts: status ? 0 : state.reconnectAttempts
    }));
}

// Increment reconnect attempts
export function incrementReconnectAttempts() {
    logsConnectionStatus.update(state => ({
        ...state,
        reconnectAttempts: state.reconnectAttempts + 1
    }));
}

// Set project ID for log streaming
export function setProjectId(projectId: string) {
    logsConnectionStatus.update(state => ({
        ...state,
        projectId
    }));
}

// Toggle auto-scroll feature
export function toggleAutoScroll(value: boolean) {
    logsConnectionStatus.update(state => ({
        ...state,
        autoScroll: value
    }));
}

// Set loading state
export function setLoading(isLoading: boolean) {
    logsConnectionStatus.update(state => ({
        ...state,
        isLoading
    }));
}

// Set error state
export function setError(error: string | null) {
    logsConnectionStatus.update(state => ({
        ...state,
        error
    }));
}

// Set total logs count
export function setTotalLogs(total: number) {
    logsConnectionStatus.update(state => ({
        ...state,
        total
    }));
}
