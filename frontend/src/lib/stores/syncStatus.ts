import { writable } from 'svelte/store';

export interface SyncStatus {
  isLoading: boolean;
  isSuccess: boolean;
  error: string | null;
}

const initialState: SyncStatus = {
  isLoading: false,
  isSuccess: false,
  error: null
};

export const syncStatus = writable<SyncStatus>(initialState);

export function resetSyncStatus() {
  syncStatus.set(initialState);
} 