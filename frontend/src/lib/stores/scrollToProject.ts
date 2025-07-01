import { writable } from 'svelte/store';

// Store to trigger scroll to specific project
export const scrollToProjectStore = writable<{ projectId: string; timestamp: number } | null>(null);

// Function to trigger scroll to project
export function triggerScrollToProject(projectId: string) {
    scrollToProjectStore.set({ projectId, timestamp: Date.now() });
}
