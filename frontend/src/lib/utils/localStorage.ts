import { browser } from '$app/environment';
import { isAuthenticated, auth } from '$lib/stores/auth';

/**
 * Get item from localStorage safely
 */
export function getLocalStorage(key: string): string | null {
	if (browser) {
		return localStorage.getItem(key);
	}
	return null;
}

/**
 * Set item to localStorage safely
 */
export function setLocalStorage(key: string, value: string) {
	if (browser) {
		localStorage.setItem(key, value);
	}
}

/**
 * Remove item from localStorage safely
 */
export function removeLocalStorage(key: string) {
	if (browser) {
		localStorage.removeItem(key);
	}
}
export function getCurrentWorkspaceId(): string | null {
	return getLocalStorage('currentWorkspaceId');
}

export function setCurrentWorkspaceId(id: string) {
	setLocalStorage('currentWorkspaceId', id);
}

// Panel width management
const DEFAULT_PROJECT_PANEL_WIDTH = 18; // rem
const DEFAULT_ROUTES_PANEL_WIDTH = 25; // percentage
const DEFAULT_REPLAY_PANEL_WIDTH = 33; // percentage

/**
 * Get saved project panel width or return default
 */
export function getProjectPanelWidth(): number {
	const saved = getLocalStorage('projectPanelWidth');
	return saved ? parseFloat(saved) : DEFAULT_PROJECT_PANEL_WIDTH;
}

/**
 * Save project panel width to localStorage
 */
export function setProjectPanelWidth(width: number) {
	setLocalStorage('projectPanelWidth', width.toString());
}

/**
 * Get saved routes panel width or return default
 */
export function getRoutesPanelWidth(): number {
	const saved = getLocalStorage('routesPanelWidth');
	return saved ? parseFloat(saved) : DEFAULT_ROUTES_PANEL_WIDTH;
}

/**
 * Save routes panel width to localStorage
 */
export function setRoutesPanelWidth(width: number) {
	setLocalStorage('routesPanelWidth', width.toString());
}

/**
 * Get saved replay panel width or return default
 */
export function getReplayPanelWidth(): number {
	const saved = getLocalStorage('replayPanelWidth');
	return saved ? parseFloat(saved) : DEFAULT_REPLAY_PANEL_WIDTH;
}

/**
 * Save replay panel width to localStorage
 */
export function setReplayPanelWidth(width: number) {
	setLocalStorage('replayPanelWidth', width.toString());
}
