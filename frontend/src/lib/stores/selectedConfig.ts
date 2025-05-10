import { writable } from 'svelte/store';
import type { ConfigResponse } from '$lib/api/mockoonApi';

export const selectedConfig = writable<ConfigResponse | null>(null);
