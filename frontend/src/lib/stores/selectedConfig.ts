import { writable } from 'svelte/store';
import type { Project } from '$lib/api/BeoApi';

export const selectedProject = writable<Project | null>(null);
