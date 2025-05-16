
export interface SystemConfig {
  id: string;
  key: string;        // Unique key for the config
  value: string;      // Value of the config
  type: string;       // string, number, boolean, json
  description: string;
  hide_value: boolean; // hide value in the UI (note the snake_case)
  created_at: string;
  updated_at: string;
}

// For backwards compatibility
export type SystemConfigItem = SystemConfig;
