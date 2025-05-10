import { writable } from 'svelte/store';
import type { Project } from '$lib/api/mockoonApi';

export const selectedProject = writable<Project | null>(null);
