import { apiClient } from './apiClient';
import { getCurrentWorkspaceId } from '$lib/utils/localStorage';

// Rule type definition matching the backend MockRule model
export type Rule = {
  id: string;
  responseId: string;
  type: string;       // "header", "query", "body"
  key: string;        // The key to match against
  operator: string;   // "equals", "contains", "regex"
  value: string;      // The value to match
  isNew?: boolean; // Optional, used for UI state management
};

// Get all rules for a response
export const getRules = async (projectId: string, endpointId: string, responseId: string): Promise<Rule[]> => {
  const workspaceId = getCurrentWorkspaceId();
  const response = await apiClient.get(
    `/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}/rules`
  );
  return response.data.data;
};

// Get a specific rule
export const getRule = async (
  projectId: string, 
  endpointId: string, 
  responseId: string, 
  ruleId: string
): Promise<Rule> => {
  const workspaceId = getCurrentWorkspaceId();
  const response = await apiClient.get(
    `/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}/rules/${ruleId}`
  );
  return response.data.data;
};

// Create a new rule
export const createRule = async (
  projectId: string,
  endpointId: string,
  responseId: string,
  rule: Omit<Rule, 'id'>
): Promise<Rule> => {
  const workspaceId = getCurrentWorkspaceId();
  const response = await apiClient.post(
    `/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}/rules`,
    rule
  );
  return response.data.data;
};

// Update an existing rule
export const updateRule = async (
  projectId: string,
  endpointId: string,
  responseId: string,
  ruleId: string,
  rule: Partial<Rule>
): Promise<Rule> => {
  const workspaceId = getCurrentWorkspaceId();
  const response = await apiClient.put(
    `/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}/rules/${ruleId}`,
    rule
  );
  return response.data.data;
};

// Delete a rule
export const deleteRule = async (
  projectId: string,
  endpointId: string,
  responseId: string,
  ruleId: string
): Promise<void> => {
  const workspaceId = getCurrentWorkspaceId();
  await apiClient.delete(
    `/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}/rules/${ruleId}`
  );
};

// Delete all rules for a response
export const deleteAllRules = async (
  projectId: string,
  endpointId: string,
  responseId: string
): Promise<void> => {
  const workspaceId = getCurrentWorkspaceId();
  await apiClient.delete(
    `/workspaces/${workspaceId}/projects/${projectId}/endpoints/${endpointId}/responses/${responseId}/rules`
  );
};
