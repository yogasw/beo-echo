import { createLogStream, getLogs, clearProjectLogsApi, type RequestLog } from '$lib/api/BeoApi';
import { currentWorkspace } from '$lib/stores/workspace';
import {
    logs,
    logsConnectionStatus,
    addLog,
    updateConnectionStatus,
    incrementReconnectAttempts,
    MAX_RECONNECT_ATTEMPTS,
    RECONNECT_DELAY_MS,
    setLoading,
    setError,
    setTotalLogs,
    clearLogs,
    addBatchLogs,
    setProjectId
} from '$lib/stores/logs';
import { get } from 'svelte/store';

// Store the EventSource instance globally so we can close/reopen it
let eventSource: EventSource | null = null;

/**
 * Initialize the logs stream for a specific project
 * @param projectId - The ID of the project to stream logs for
 * @param pageSize - Number of logs to fetch initially
 */
export function initializeLogsStream(projectId: string, pageSize: number = 100, clear = false) {
    // Close any existing connection and clear previous logs
    closeLogStream();
    if (clear) {
        clearLogs();
    }
    // Load initial logs and then setup stream
    loadInitialLogs(projectId, pageSize).then(() => {
        setupLogStream(projectId, pageSize);
    });
}

/**
 * Load the initial set of logs for a project
 */
async function loadInitialLogs(projectId: string, pageSize: number = 100) {
    const workspaceData = get(currentWorkspace);

    if (!workspaceData) {
        const error = 'No workspace selected';
        console.error(error);
        setError(error);
        return;
    }

    try {
        setProjectId(projectId)
        setLoading(true);
        const result = await getLogs(1, pageSize, projectId);
        let currentLogs = get(logs);
        if (result.logs) {
            // Check if logs already exist to prevent duplicates
            const newLogs = result.logs.filter((log: RequestLog) => !currentLogs.some((existingLog) => existingLog.id === log.id));
            currentLogs = [...newLogs, ...currentLogs]
        }

        addBatchLogs(currentLogs);
        setLoading(false);
    } catch (err) {
        console.error('Failed to load logs:', err);
        setError('Failed to load logs: ' + (err instanceof Error ? err.message : String(err)));
        setLoading(false);
    }
}

/**
 * Setup the EventSource connection for log streaming
 */
function setupLogStream(projectId: string, pageSize: number = 100) {
    // Close any existing connection
    if (eventSource) {
        eventSource.close();
    }

    const workspaceData = get(currentWorkspace);

    if (!workspaceData) {
        console.error('Cannot setup log stream: No workspace selected');
        return;
    }

    console.log('Setting up log stream for project:', projectId, 'in workspace:', workspaceData.id);

    // Create new connection
    eventSource = createLogStream(projectId, pageSize);

    // Setup event handlers
    eventSource.addEventListener('log', (event) => {
        try {
            console.log('Log event received:', event.data);
            const newLog = JSON.parse(event.data);

            // Add the new log to the store
            addLog(newLog);

            // Auto-scroll to top if enabled
            if (get(logsConnectionStatus).autoScroll) {
                window.scrollTo(0, 0);
            }
        } catch (err) {
            console.error('Error processing log event:', err, event.data);
        }
    });

    // Direct message event (fallback)
    eventSource.onmessage = (event) => {
        console.log('Generic message received:', event.data);
        try {
            const newLog = JSON.parse(event.data);
            if (newLog && newLog.id) {
                addLog(newLog);
            }
        } catch (err) {
            console.error('Error processing generic message:', err);
        }
    };

    eventSource.addEventListener('ping', (event) => {
        // Keep connection alive, no action needed
        console.log('Ping received from server:', event.data);
        updateConnectionStatus(true);
    });

    eventSource.onopen = () => {
        console.log('Log stream connection established');
        updateConnectionStatus(true);
    };

    eventSource.onerror = (err) => {
        console.error('EventSource error:', err);
        updateConnectionStatus(false);

        // Implement smart reconnection strategy with backoff
        const { reconnectAttempts } = get(logsConnectionStatus);

        if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
            incrementReconnectAttempts();
            const delay = RECONNECT_DELAY_MS * (reconnectAttempts + 1); // Increase delay with each attempt

            console.log(
                `Attempting to reconnect log stream (attempt ${reconnectAttempts + 1}/${MAX_RECONNECT_ATTEMPTS}) in ${delay}ms...`
            );

            setTimeout(() => {
                setupLogStream(projectId, pageSize);
            }, delay);
        } else {
            console.error('Max reconnection attempts reached. Please refresh manually.');
        }
    };
}

/**
 * Manually reconnect the log stream
 */
export function reconnectLogStream() {
    const { projectId } = get(logsConnectionStatus);

    if (projectId) {
        setupLogStream(projectId);
    } else {
        console.error('Cannot reconnect: No project ID available');
    }
}

/**
 * Close and clean up the log stream connection
 */
export function closeLogStream() {
    if (eventSource) {
        eventSource.close();
        eventSource = null;
    }
    updateConnectionStatus(false);
}

/**
 * Refresh logs manually
 */
export function refreshLogs() {
    const { projectId } = get(logsConnectionStatus);

    if (!projectId) {
        console.error('Cannot refresh logs: No project ID available');
        return;
    }

    loadInitialLogs(projectId);

    // Also try to reconnect if disconnected
    if (!get(logsConnectionStatus).isConnected) {
        setupLogStream(projectId);
    }
}

/**
 * Clear all non-bookmarked logs for the current project
 * @returns Promise that resolves to the number of logs cleared
 */
export async function clearProjectLogs(): Promise<number> {
    const { projectId } = get(logsConnectionStatus);
    if (!projectId) {
        console.error('Cannot clear logs: No project ID available');
        throw new Error('No project ID available');
    }

    try {
        clearLogs()
        const rowsDeleted = await clearProjectLogsApi(projectId);
        refreshLogs();

        return rowsDeleted;
    } catch (err) {
        console.error('Failed to clear logs:', err);
        throw err;
    }
}
