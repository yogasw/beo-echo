/**
 * Badge styling utilities for consistent colors across the application
 */

// HTTP Method color mapping
export const HTTP_METHOD_COLORS = {
  GET: { color: 'text-green-600', bgColor: 'bg-green-600' },
  POST: { color: 'text-blue-600', bgColor: 'bg-blue-600' },
  PUT: { color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
  DELETE: { color: 'text-red-600', bgColor: 'bg-red-600' },
  PATCH: { color: 'text-purple-600', bgColor: 'bg-purple-600' },
  OPTIONS: { color: 'text-gray-600', bgColor: 'bg-gray-600' },
  HEAD: { color: 'text-gray-600', bgColor: 'bg-gray-600' },
  CONNECT: { color: 'text-gray-600', bgColor: 'bg-gray-600' },
  TRACE: { color: 'text-gray-600', bgColor: 'bg-gray-600' },
} as const;

// Status code color mapping by category
export const STATUS_CODE_COLORS = {
  SUCCESS: { color: 'text-green-600', bgColor: 'bg-green-600' },
  REDIRECTION: { color: 'text-blue-600', bgColor: 'bg-blue-600' },
  CLIENT_ERROR: { color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
  SERVER_ERROR: { color: 'text-red-600', bgColor: 'bg-red-600' },
  CUSTOM: { color: 'text-gray-600', bgColor: 'bg-gray-600' },
} as const;

/**
 * Get the appropriate color for an HTTP method
 * @param method - The HTTP method (case insensitive)
 * @returns Object containing text and background color classes
 */
export function getHttpMethodColor(method: string): { color: string; bgColor: string } {
  const upperMethod = method.toUpperCase() as keyof typeof HTTP_METHOD_COLORS;
  return HTTP_METHOD_COLORS[upperMethod] || HTTP_METHOD_COLORS.OPTIONS;
}

/**
 * Get the appropriate color for a status code based on its range
 * @param statusCode - The HTTP status code
 * @returns Object containing text and background color classes and category
 */
export function getStatusCodeColor(statusCode: number): { 
  color: string; 
  bgColor: string; 
  category: string 
} {
  if (statusCode >= 200 && statusCode < 300) {
    return { ...STATUS_CODE_COLORS.SUCCESS, category: 'Success' };
  } else if (statusCode >= 300 && statusCode < 400) {
    return { ...STATUS_CODE_COLORS.REDIRECTION, category: 'Redirection' };
  } else if (statusCode >= 400 && statusCode < 500) {
    return { ...STATUS_CODE_COLORS.CLIENT_ERROR, category: 'Client Error' };
  } else if (statusCode >= 500 && statusCode < 600) {
    return { ...STATUS_CODE_COLORS.SERVER_ERROR, category: 'Server Error' };
  } else {
    return { ...STATUS_CODE_COLORS.CUSTOM, category: 'Custom' };
  }
}

/**
 * Get method description for display purposes
 * @param method - The HTTP method
 * @returns Human-readable description of the method
 */
export function getHttpMethodDescription(method: string): string {
  const descriptions: Record<string, string> = {
    GET: 'Retrieve data from server',
    POST: 'Create new resource',
    PUT: 'Update existing resource',
    DELETE: 'Remove resource',
    PATCH: 'Partially update resource',
    OPTIONS: 'Get allowed methods',
    HEAD: 'Get headers only',
    CONNECT: 'Establish tunnel',
    TRACE: 'Perform message loop-back test',
  };

  return descriptions[method.toUpperCase()] || 'Custom HTTP method';
}

/**
 * Get status code description and name
 * @param statusCode - The HTTP status code
 * @returns Object with name and description, or null if not found
 */
export function getStatusCodeInfo(statusCode: number): { 
  name: string; 
  description: string; 
  category: string 
} | null {
  const statusCodes: Record<number, { name: string; description: string; category: string }> = {
    // 2xx Success
    200: { name: 'OK', description: 'Request successful', category: 'Success' },
    201: { name: 'Created', description: 'Resource created successfully', category: 'Success' },
    202: { name: 'Accepted', description: 'Request accepted for processing', category: 'Success' },
    204: { name: 'No Content', description: 'Success with no response body', category: 'Success' },
    206: { name: 'Partial Content', description: 'Partial content served', category: 'Success' },
    
    // 3xx Redirection  
    301: { name: 'Moved Permanently', description: 'Resource moved permanently', category: 'Redirection' },
    302: { name: 'Found', description: 'Resource temporarily moved', category: 'Redirection' },
    304: { name: 'Not Modified', description: 'Resource not modified', category: 'Redirection' },
    307: { name: 'Temporary Redirect', description: 'Temporary redirect', category: 'Redirection' },
    308: { name: 'Permanent Redirect', description: 'Permanent redirect', category: 'Redirection' },
    
    // 4xx Client Error
    400: { name: 'Bad Request', description: 'Invalid request syntax', category: 'Client Error' },
    401: { name: 'Unauthorized', description: 'Authentication required', category: 'Client Error' },
    403: { name: 'Forbidden', description: 'Access denied', category: 'Client Error' },
    404: { name: 'Not Found', description: 'Resource not found', category: 'Client Error' },
    405: { name: 'Method Not Allowed', description: 'HTTP method not allowed', category: 'Client Error' },
    406: { name: 'Not Acceptable', description: 'Response format not accepted', category: 'Client Error' },
    408: { name: 'Request Timeout', description: 'Request timeout', category: 'Client Error' },
    409: { name: 'Conflict', description: 'Request conflicts with current state', category: 'Client Error' },
    410: { name: 'Gone', description: 'Resource no longer available', category: 'Client Error' },
    412: { name: 'Precondition Failed', description: 'Precondition failed', category: 'Client Error' },
    413: { name: 'Payload Too Large', description: 'Request payload too large', category: 'Client Error' },
    415: { name: 'Unsupported Media Type', description: 'Media type not supported', category: 'Client Error' },
    422: { name: 'Unprocessable Entity', description: 'Validation failed', category: 'Client Error' },
    429: { name: 'Too Many Requests', description: 'Rate limit exceeded', category: 'Client Error' },
    
    // 5xx Server Error
    500: { name: 'Internal Server Error', description: 'Server encountered an error', category: 'Server Error' },
    501: { name: 'Not Implemented', description: 'Server does not support functionality', category: 'Server Error' },
    502: { name: 'Bad Gateway', description: 'Invalid response from upstream', category: 'Server Error' },
    503: { name: 'Service Unavailable', description: 'Server temporarily unavailable', category: 'Server Error' },
    504: { name: 'Gateway Timeout', description: 'Upstream server timeout', category: 'Server Error' },
    505: { name: 'HTTP Version Not Supported', description: 'HTTP version not supported', category: 'Server Error' }
  };

  return statusCodes[statusCode] || null;
}

/**
 * Create badge size classes for consistent sizing
 * @param size - The size variant
 * @returns CSS classes for the specified size
 */
export function getBadgeSizeClasses(size: 'xs' | 'sm' | 'md' | 'lg'): string {
  const sizeMap = {
    xs: 'px-1.5 py-0.5 text-xs',
    sm: 'px-2 py-1 text-xs',
    md: 'px-2.5 py-1.5 text-sm',
    lg: 'px-3 py-2 text-base'
  };
  
  return sizeMap[size];
}

/**
 * Create base badge classes for consistent styling
 * @returns CSS classes for badge base styling
 */
export function getBaseBadgeClasses(): string {
  return 'rounded font-medium text-white inline-flex items-center';
}
