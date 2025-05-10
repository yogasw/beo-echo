import { writable } from 'svelte/store';

export const isFirstOpenPage = writable<boolean>(true);
