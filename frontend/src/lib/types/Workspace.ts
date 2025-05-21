import type { UserWorkspace } from './User';

/**
 * Workspace representation
 */
export interface Workspace {
    /** Unique identifier for workspace */
    id: string;
    
    /** Name of the workspace */
    name: string;
    
    /** List of projects in this workspace */
    projects?: Project[];
    
    /** List of user memberships in this workspace */
    members?: UserWorkspace[];
    
    /** Whether auto-invite is enabled for this workspace */
    auto_invite_enabled: boolean;
    
    /** Comma-separated list of email domains for auto-invite */
    auto_invite_domains?: string;
    
    /** Role assigned to auto-invited users */
    auto_invite_role?: string;
    
    /** Creation timestamp */
    created_at: string;
    
    /** Last update timestamp */
    updated_at: string;
}

/**
 * Project representation in a workspace
 */
export interface Project {
    /** Unique identifier for project */
    id: string;
    
    /** Name of the project */
    name: string;
    
    /** ID of the workspace this project belongs to */
    workspace_id: string;
    
    /** Project operating mode */
    mode: 'mock' | 'proxy' | 'forwarder' | 'disabled';
    
    /** Project status */
    status: 'running' | 'stopped' | 'error';
    
    /** Project alias/subdomain */
    alias: string;
    
    /** Project URL */
    url?: string;
    
    /** Creation timestamp */
    created_at: string;
    
    /** Last update timestamp */
    updated_at: string;
}
