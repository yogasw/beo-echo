import { writable } from 'svelte/store';

export interface GoogleOAuthConfig {
    client_id: string;
    client_secret: string;
    allow_domains: string[];
    instructions: string;
}

interface SSOState {
    googleConfig: GoogleOAuthConfig | null;
    isGoogleEnabled: boolean;
    isConfigModalOpen: boolean;
    autoCreateAccounts: boolean;
    allowLocalAuth: boolean;
    defaultSignInMethod: 'local' | 'google' | 'github' | 'microsoft' | 'saml';
    loading: boolean;
    error: string | null;
}

const initialState: SSOState = {
    googleConfig: null,
    isGoogleEnabled: false,
    isConfigModalOpen: false,
    autoCreateAccounts: true,
    allowLocalAuth: true,
    defaultSignInMethod: 'local',
    loading: false,
    error: null,
};

function createSSOStore() {
    const { subscribe, set, update } = writable<SSOState>(initialState);

    return {
        subscribe,
        setGoogleConfig: (config: GoogleOAuthConfig) => update(state => ({ ...state, googleConfig: config })),
        setGoogleEnabled: (enabled: boolean) => update(state => ({ ...state, isGoogleEnabled: enabled })),
        toggleConfigModal: () => update(state => ({ ...state, isConfigModalOpen: !state.isConfigModalOpen })),
        setAutoCreateAccounts: (enabled: boolean) => update(state => ({ ...state, autoCreateAccounts: enabled })),
        setAllowLocalAuth: (enabled: boolean) => update(state => ({ ...state, allowLocalAuth: enabled })),
        setDefaultSignInMethod: (method: SSOState['defaultSignInMethod']) => 
            update(state => ({ ...state, defaultSignInMethod: method })),
        setLoading: (loading: boolean) => update(state => ({ ...state, loading })),
        setError: (error: string | null) => update(state => ({ ...state, error })),
        reset: () => set(initialState)
    };
}

export const ssoStore = createSSOStore();
