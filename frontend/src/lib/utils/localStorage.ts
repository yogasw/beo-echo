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
