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
    show: (message: string, type: ToastType = 'info', duration: number = 5000) => {
      set({ message, type, duration });
    },
    success: (message: string, duration?: number) => {
      set({ message, type: 'success', duration });
    },
    error: (message: string | Error | any, duration?: number) => {
      let errorMessage: string;

      if (typeof message === 'string') {
        // If message is already a string, use it directly
        errorMessage = message;
      } else if (message && message.response && message.response.data && message.response.data.message) {
        // Handle Axios error format
        errorMessage = message.response.data.message;
      } else if (message instanceof Error) {
        // If message is an Error object, get its message property
        errorMessage = message.message;
      } else {
        // Fallback for any other type
        errorMessage = 'An unknown error occurred';
      }

      set({ message: errorMessage, type: 'error', duration });
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