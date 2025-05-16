import { writable } from 'svelte/store';
import { browser } from '$app/environment';

interface FeatureToggles {
    showPasswordRequirements: boolean;
    // Add more feature toggles here as needed
    // Examples:
    // allowApiKeyGeneration: boolean;
    // enableNotifications: boolean;
    // allowProjectDeletion: boolean;
    [key: string]: boolean; // Allow dynamic feature flags
}

export enum FeatureFlags {
    // Both uppercase and lowercase versions are supported by the backend
    FEATURE_EMAIL_UPDATES_ENABLED = 'FEATURE_EMAIL_UPDATES_ENABLED',
    FEATURE_SHOW_PASSWORD_REQUIREMENTS = 'FEATURE_SHOW_PASSWORD_REQUIREMENTS',
    FEATURE_REGISTER_EMAIL_ENABLED = 'FEATURE_REGISTER_EMAIL_ENABLED',
}

// Default feature toggle settings
const defaultToggles: FeatureToggles = {
    showPasswordRequirements: false,
};

// Try to get stored settings from localStorage, or use defaults
function getInitialToggles(): FeatureToggles {
    if (browser) {
        // Clear localStorage to ensure new default settings are applied
        // This line can be removed after first deployment when everyone has the new settings
        localStorage.removeItem('featureToggles');

        const storedToggles = localStorage.getItem('featureToggles');
        if (storedToggles) {
            try {
                return { ...defaultToggles, ...JSON.parse(storedToggles) };
            } catch (e) {
                console.error('Failed to parse feature toggles from localStorage:', e);
            }
        }
    }
    return { ...defaultToggles };
}

// Create the writable store with initial values
export const featureToggles = writable<FeatureToggles>(getInitialToggles());

// Subscribe to changes and update localStorage
if (browser) {
    featureToggles.subscribe(value => {
        localStorage.setItem('featureToggles', JSON.stringify(value));
    });
}

// Helper function to update a single toggle
// This will be triggered from the backend
export function updateFeatureToggle(feature: FeatureFlags, enabled: boolean): void {
    featureToggles.update(toggles => ({
        ...toggles,
        [feature]: enabled
    }));
}

// get the feature toggle value
export function getFeatureToggle(feature: FeatureFlags): boolean {
    let value = false;
    featureToggles.subscribe(toggles => {
        value = toggles[feature];
    })();
    return value;
}

/**
 * Synchronizes feature flags from the backend with the client-side store
 * @param flags Object containing feature flags from the backend API
 */
export function syncFeatureFlags(flags: Record<string, boolean>): void {
    if (!flags || typeof flags !== 'object') {
        console.warn('Invalid feature flags received:', flags);
        return;
    }
    // Update each flag in our store
    Object.entries(flags).forEach(([key, value]) => {
        updateFeatureToggle(key as FeatureFlags, value);
    });
}


export default featureToggles;
