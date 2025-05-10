import { writable } from 'svelte/store';
import type { ProjectResponse } from '$lib/api/mockoonApi';

export const selectedConfig = writable<ProjectResponse | null>(null);
