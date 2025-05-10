import { writable } from 'svelte/store';
import type { Config } from '$lib/api/mockoonApi';
import { getConfigs } from '$lib/api/mockoonApi';

export const configurations = writable<Config[]>([]);

export async function fetchConfigsStore() {
  configurations.set(await getConfigs());
} 