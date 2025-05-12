import { writable } from 'svelte/store';
import type { Project } from '$lib/api/BeoApi';
import { getProjects } from '$lib/api/BeoApi';

export const projects = writable<Project[]>([]);

export async function fetchConfigsStore() {
  projects.set(await getProjects());
} 