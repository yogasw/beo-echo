import { writable } from 'svelte/store';
import type { Project } from '$lib/api/mockoonApi';
import { getProjects } from '$lib/api/mockoonApi';

export const projects = writable<Project[]>([]);

export async function fetchConfigsStore() {
  projects.set(await getProjects());
} 