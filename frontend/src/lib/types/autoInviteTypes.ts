import type { WorkspaceRole } from './User';

/**
 * Auto-invite configuration for a workspace
 * This configuration allows automatically adding users to workspaces based on their email domains
 */
export interface AutoInviteConfig {
    /** Whether auto-invite feature is enabled for this workspace */
    enabled: boolean;
    
    /** List of email domains that will be automatically added to this workspace */
    domains: string[];
    
    /** The role that will be assigned to auto-invited users */
    role: WorkspaceRole;
    
    /** The workspace ID this configuration belongs to */
    workspace_id: string;
    
    /** The workspace name (for display purposes) */
    workspace_name: string;
}
