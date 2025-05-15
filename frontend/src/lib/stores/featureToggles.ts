import { writable } from 'svelte/store';
import { browser } from '$app/environment';

interface FeatureToggles {
    showPasswordRequirements: boolean;
    // Add more feature toggles here as needed
    // Examples:
    // allowApiKeyGeneration: boolean;
    // enableNotifications: boolean;
    // allowProjectDeletion: boolean;
}

// Default feature toggle settings
const defaultToggles: FeatureToggles = {
    showPasswordRequirements: false // Enabled by default now
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
export function updateFeatureToggle(feature: keyof FeatureToggles, enabled: boolean): void {
    featureToggles.update(toggles => ({
        ...toggles,
        [feature]: enabled
    }));
}

export { featureToggles as default };
