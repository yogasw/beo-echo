import { apiClient } from './apiClient';

/**
 * Personal Access Tokens (PAT) — long-lived tokens used to authenticate the
 * MCP server / CLIs against the Beo Echo API. The plaintext token is only ever
 * returned once, at creation time.
 */

export interface ApiToken {
	id: string;
	name: string;
	prefix: string;
	source: string; // "pat" | "oauth"
	client_id: string;
	last_used_at: string | null;
	expires_at: string | null;
	created_at: string;
}

export interface CreatedApiToken {
	token: string; // plaintext — shown once
	id: string;
	name: string;
	prefix: string;
	expires_at: string | null;
	created_at: string;
}

/** List the current user's access tokens (metadata only, never the secret). */
export async function listApiTokens(): Promise<ApiToken[]> {
	const res = await apiClient.get('/users/me/tokens');
	return res.data?.data ?? [];
}

/**
 * Create a new access token. expiresDays of 0 (the default) means it never
 * expires. Returns the plaintext token — copy it immediately.
 */
export async function createApiToken(name: string, expiresDays = 0): Promise<CreatedApiToken> {
	const res = await apiClient.post('/users/me/tokens', {
		name,
		expires_days: expiresDays
	});
	return res.data.data;
}

/** Revoke (delete) an access token by id. */
export async function revokeApiToken(tokenId: string): Promise<void> {
	await apiClient.delete(`/users/me/tokens/${tokenId}`);
}
