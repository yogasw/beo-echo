import { writable } from 'svelte/store';
import type { Project } from '$lib/api/BeoApi';

export const isLoadingContentArea = writable<boolean>(false);
