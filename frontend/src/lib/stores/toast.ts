import { writable } from 'svelte/store';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface Toast {
  message: string;
  type: ToastType;
  duration?: number;
}

const createToastStore = () => {
  const { subscribe, set } = writable<Toast | null>(null);

  return {
    subscribe,
    show: (message: string, type: ToastType = 'info', duration: number = 4000) => {
      set({ message, type, duration });
    },
    success: (message: string, duration?: number) => {
      set({ message, type: 'success', duration });
    },
    error: (message: string, duration?: number) => {
      set({ message, type: 'error', duration });
    },
    warning: (message: string, duration?: number) => {
      set({ message, type: 'warning', duration });
    },
    info: (message: string, duration?: number) => {
      set({ message, type: 'info', duration });
    },
    clear: () => {
      set(null);
    }
  };
};

export const toast = createToastStore(); 