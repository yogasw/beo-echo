/**
 * Utility function for project alias sanitization
 * Used across the application to ensure consistent alias handling
 */

/**
 * Sanitize alias input to enforce format restrictions
 * - Replaces spaces with hyphens
 * - Converts to lowercase
 * - Removes invalid characters (only allows a-z, 0-9, -)
 * - Removes consecutive hyphens
 * - Removes leading/trailing hyphens
 */
export function sanitizeAlias(input: string): string {
	// Replace spaces with hyphens
	let sanitized = input.replace(/\s+/g, '-');

	// Convert to lowercase
	sanitized = sanitized.toLowerCase();

	// Remove all characters that are not lowercase letters, numbers, or hyphens
	sanitized = sanitized.replace(/[^a-z0-9-]/g, '');

	// Remove consecutive hyphens
	sanitized = sanitized.replace(/-+/g, '-');

	return sanitized;
}
