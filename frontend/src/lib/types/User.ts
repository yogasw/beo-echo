export interface User {
  id: string;
  email: string;
  name: string;
  isOwner: boolean;
  isEnabled?: boolean; // Flag to determine if user account is active or disabled
  feature_flags?: Record<string, boolean>; // Feature flags from the backend
}

export interface Workspace {
  id: string;
  name: string;
  role?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface WorkspaceMember {
  id: string;
  userId: string;
  workspaceId: string;
  role: string;
  createdAt?: string;
  updatedAt?: string;
  user?: User;
}

export type WorkspaceRole = 'admin' | 'member';
