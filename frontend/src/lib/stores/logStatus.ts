import { writable } from 'svelte/store';

interface LogStatus {
  isConnected: boolean; // true if connected to logs stream
  unreadCount: number; // count of unread logs
}

// Create the store with initial values
const createLogStatusStore = () => {
  const { subscribe, update, set } = writable<LogStatus>({
    isConnected: false,
    unreadCount: 0
  });

  return {
    subscribe,
    
    // Set connection status
    setConnectionStatus: (status: boolean) => update(state => ({ 
      ...state, 
      isConnected: status 
    })),
    
    // Increment unread count
    incrementUnread: () => update(state => ({ 
      ...state, 
      unreadCount: state.unreadCount + 1 
    })),
    
    // Reset unread count (when logs tab is viewed)
    resetUnread: () => update(state => ({ 
      ...state, 
      unreadCount: 0 
    })),
    
    // Reset the entire store
    reset: () => set({ 
      isConnected: false, 
      unreadCount: 0 
    })
  };
};

export const logStatus = createLogStatusStore();
