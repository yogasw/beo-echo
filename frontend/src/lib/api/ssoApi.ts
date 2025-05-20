import type { GoogleOAuthConfig } from '$lib/stores/ssoStore';
import { HttpApi } from './BeoApi';

export const ssoApi = {
    // Google OAuth
    getGoogleConfig: async () => {
        try {
            const { data } = await HttpApi.get('/oauth/google/config');
            return data.data;
        } catch (error) {
            throw new Error('Failed to fetch Google OAuth config');
        }
    },

    updateGoogleConfig: async (config: GoogleOAuthConfig) => {
        try {
            const { data } = await HttpApi.put('/oauth/google/config', config);
            return data;
        } catch (error) {
            throw new Error('Failed to update Google OAuth config');
        }
    },

    updateGoogleState: async (enabled: boolean) => {
        try {
            const { data } = await HttpApi.put('/oauth/google/state', { enabled });
            return data;
        } catch (error) {
            throw new Error('Failed to update Google OAuth state');
        }
    },

    // System Config
    updateSystemConfig: async (key: string, value: any) => {
        try {
            const { data } = await HttpApi.put(`/system-config/${key}`, { value });
            return data;
        } catch (error) {
            throw new Error('Failed to update system config');
        }
    },
};
