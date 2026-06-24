import { writable } from 'svelte/store';

export const activeTab = writable('routes');

// When navigating to the Settings tab, this signals which collapsible section
// should be expanded (e.g. 'profile', 'password', 'mcp'). Null = leave defaults.
export const settingsSection = writable<string | null>(null);