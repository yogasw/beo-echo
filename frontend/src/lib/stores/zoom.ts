// src/lib/stores/zoom.ts
import { browser } from '$app/environment';
import { writable } from 'svelte/store';

// Default zoom level is 0.9 depending on user preference
const storedZoom = browser && localStorage.getItem('app-zoom');
const initialZoom = storedZoom ? parseFloat(storedZoom) : 0.9;

export const zoomLevel = writable<number>(initialZoom);

zoomLevel.subscribe((value) => {
	if (browser) {
		localStorage.setItem('app-zoom', value.toString());
	}
});
