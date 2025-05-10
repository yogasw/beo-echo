import { writable } from 'svelte/store';
import type { Config } from '$lib/api/mockoonApi';
import { getProjects } from '$lib/api/mockoonApi';

export const configurations = writable<Config[]>([]);

export async function fetchConfigsStore() {
  configurations.set(await getProjects());
} 