
import type { SystemConfig } from '$lib/types/SystemConfig';
import { BASE_URL_API } from '$lib/utils/authUtils';
import { auth } from '$lib/stores/auth';

/**
 * API service for system configuration
 */
export const SystemConfigAPI = {
  /**
   * Get a system config by key
   * @param key Configuration key
   * @returns Configuration value
   */
  async getByKey(key: string): Promise<SystemConfig> {
    const token = auth.getToken();
    
    const response = await fetch(`${BASE_URL_API}/system-config/${key}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`Failed to get system config: ${response.statusText}`);
    }
    
    return await response.json();
  },

  /**
   * Update a system config
   * @param key Configuration key
   * @param value New value
   * @returns Updated configuration
   */
  async update(key: string, value: string): Promise<SystemConfig> {
    const token = auth.getToken();
    
    const response = await fetch(`${BASE_URL_API}/system-config/${key}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ value })
    });
    
    if (!response.ok) {
      throw new Error(`Failed to update system config: ${response.statusText}`);
    }
    
    return await response.json();
  },
  
  /**
   * Get the value of a feature flag
   * @param key Feature flag key
   * @param defaultValue Default value if flag doesn't exist
   * @returns Boolean value of the feature flag
   */
  async getFeatureFlag(key: string, defaultValue = false): Promise<boolean> {
    try {
      const config = await this.getByKey(key);
      return config.value === 'true';
    } catch (error) {
      console.warn(`Feature flag ${key} not found, using default: ${defaultValue}`);
      return defaultValue;
    }
  }
};
