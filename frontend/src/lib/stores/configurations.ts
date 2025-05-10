import { writable } from 'svelte/store';
import type { ProjectResponse } from '$lib/api/mockoonApi';
import { getProjects } from '$lib/api/mockoonApi';

export const configurations = writable<ProjectResponse[]>([]);

export async function fetchConfigsStore() {
  configurations.set(await getProjects());
} 