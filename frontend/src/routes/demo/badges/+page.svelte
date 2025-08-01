<!-- Demo page for reusable badge components -->
<script lang="ts">
  import StatusCodeBadge from '$lib/components/common/StatusCodeBadge.svelte';
  import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';
  import EndpointListItem from '$lib/components/common/EndpointListItem.svelte';
  import EndpointList from '$lib/components/common/EndpointList.svelte';

  // Sample data for demonstration
  const sampleStatusCodes = [200, 201, 204, 301, 302, 400, 401, 403, 404, 422, 500, 502, 503, 418];
  const sampleMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS', 'HEAD', 'CONNECT'];
  const sizes = ['xs', 'sm', 'md', 'lg'] as const;

  // Mock API request data for realistic examples
  const mockRequests = [
    { method: 'GET', path: '/api/users', statusCode: 200, description: 'Fetch all users' },
    { method: 'POST', path: '/api/users', statusCode: 201, description: 'Create new user' },
    { method: 'PUT', path: '/api/users/123', statusCode: 200, description: 'Update user' },
    { method: 'DELETE', path: '/api/users/123', statusCode: 204, description: 'Delete user' },
    { method: 'GET', path: '/api/posts', statusCode: 404, description: 'Posts not found' },
    { method: 'POST', path: '/api/auth/login', statusCode: 401, description: 'Invalid credentials' },
    { method: 'PATCH', path: '/api/settings', statusCode: 422, description: 'Validation failed' },
    { method: 'GET', path: '/api/health', statusCode: 500, description: 'Server error' }
  ];

  // Mock endpoints data for endpoint list demonstration
  const mockEndpoints = [
    {
      id: '1',
      method: 'GET',
      path: '/api/users',
      name: 'Get Users',
      description: 'Retrieve a list of all users',
      status: 'active' as const,
      statusCode: 200,
      responseTime: 45,
      lastUsed: new Date(Date.now() - 300000), // 5 minutes ago
      isDefault: false
    },
    {
      id: '2', 
      method: 'POST',
      path: '/api/users',
      name: 'Create User',
      description: 'Create a new user account',
      status: 'active' as const,
      statusCode: 201,
      responseTime: 120,
      lastUsed: new Date(Date.now() - 1800000), // 30 minutes ago
      isDefault: true
    },
    {
      id: '3',
      method: 'PUT',
      path: '/api/users/{id}',
      name: 'Update User',
      description: 'Update an existing user',
      status: 'inactive' as const,
      statusCode: 200,
      responseTime: 85,
      lastUsed: new Date(Date.now() - 86400000), // 1 day ago
      isDefault: false
    },
    {
      id: '4',
      method: 'DELETE',
      path: '/api/users/{id}',
      name: 'Delete User',
      description: 'Remove a user from the system',
      status: 'error' as const,
      statusCode: 500,
      responseTime: undefined,
      lastUsed: new Date(Date.now() - 259200000), // 3 days ago
      isDefault: false
    },
    {
      id: '5',
      method: 'PATCH',
      path: '/api/users/{id}/settings',
      name: 'Update Settings',
      description: 'Update user preferences and settings',
      status: 'active' as const,
      statusCode: 200,
      responseTime: 32,
      lastUsed: new Date(Date.now() - 600000), // 10 minutes ago
      isDefault: false
    }
  ];

  function handleEndpointAction(event: CustomEvent) {
    console.log('Endpoint action:', event.type, event.detail);
  }
</script>

<svelte:head>
  <title>Badge Components Demo - Beo Echo</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
  <div class="max-w-6xl mx-auto px-4">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
        <i class="fas fa-tags mr-3"></i>Reusable Badge Components
      </h1>
      <p class="text-gray-600 dark:text-gray-400">
        Consistent badge styling for HTTP methods and status codes throughout the application
      </p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Status Code Badges -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          <i class="fas fa-hashtag mr-2"></i>Status Code Badges
        </h2>
        
        <div class="space-y-6">
          <!-- Different sizes -->
          <div>
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">Different Sizes</h3>
            <div class="flex flex-wrap items-center gap-3">
              {#each sizes as size}
                <div class="flex flex-col items-center gap-1">
                  <StatusCodeBadge statusCode={200} {size} />
                  <span class="text-xs text-gray-500">{size}</span>
                </div>
              {/each}
            </div>
          </div>

          <!-- All status codes -->
          <div>
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">All Status Codes</h3>
            <div class="flex flex-wrap gap-2">
              {#each sampleStatusCodes as statusCode}
                <StatusCodeBadge {statusCode} />
              {/each}
            </div>
          </div>

          <!-- With descriptions -->
          <div>
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">With Descriptions</h3>
            <div class="space-y-2">
              <StatusCodeBadge statusCode={200} showDescription={true} />
              <StatusCodeBadge statusCode={404} showDescription={true} />
              <StatusCodeBadge statusCode={500} showDescription={true} />
              <StatusCodeBadge statusCode={418} showDescription={true} />
            </div>
          </div>
        </div>
      </div>

      <!-- HTTP Method Badges -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          <i class="fas fa-code mr-2"></i>HTTP Method Badges
        </h2>
        
        <div class="space-y-6">
          <!-- Different sizes -->
          <div>
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">Different Sizes</h3>
            <div class="flex flex-wrap items-center gap-3">
              {#each sizes as size}
                <div class="flex flex-col items-center gap-1">
                  <HttpMethodBadge method="GET" {size} />
                  <span class="text-xs text-gray-500">{size}</span>
                </div>
              {/each}
            </div>
          </div>

          <!-- All methods -->
          <div>
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">All HTTP Methods</h3>
            <div class="flex flex-wrap gap-2">
              {#each sampleMethods as method}
                <HttpMethodBadge {method} />
              {/each}
            </div>
          </div>

          <!-- With descriptions -->
          <div>
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">With Descriptions</h3>
            <div class="space-y-2">
              <HttpMethodBadge method="GET" showDescription={true} />
              <HttpMethodBadge method="POST" showDescription={true} />
              <HttpMethodBadge method="DELETE" showDescription={true} />
              <HttpMethodBadge method="CUSTOM" showDescription={true} />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Realistic Usage Examples -->
    <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
        <i class="fas fa-list mr-2"></i>Realistic Usage Examples
      </h2>
      
      <div class="space-y-4">
        <div class="text-sm text-gray-600 dark:text-gray-400 mb-4">
          Mock API requests showing how badges would appear in real application scenarios:
        </div>
        
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-gray-600">
                <th class="text-left py-2 text-gray-700 dark:text-gray-300 font-medium">Method</th>
                <th class="text-left py-2 text-gray-700 dark:text-gray-300 font-medium">Endpoint</th>
                <th class="text-left py-2 text-gray-700 dark:text-gray-300 font-medium">Status</th>
                <th class="text-left py-2 text-gray-700 dark:text-gray-300 font-medium">Description</th>
              </tr>
            </thead>
            <tbody>
              {#each mockRequests as request}
                <tr class="border-b border-gray-100 dark:border-gray-700">
                  <td class="py-3">
                    <HttpMethodBadge method={request.method} size="sm" />
                  </td>
                  <td class="py-3 text-gray-900 dark:text-white font-mono text-sm">
                    {request.path}
                  </td>
                  <td class="py-3">
                    <StatusCodeBadge statusCode={request.statusCode} size="sm" />
                  </td>
                  <td class="py-3 text-gray-600 dark:text-gray-400">
                    {request.description}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Usage Instructions -->
    <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
        <i class="fas fa-code mr-2"></i>Usage Instructions
      </h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h3 class="font-medium text-gray-900 dark:text-white mb-3">StatusCodeBadge</h3>
          <pre class="bg-gray-100 dark:bg-gray-700 p-3 rounded text-xs overflow-x-auto"><code>{`<StatusCodeBadge 
  statusCode={403} 
  size="sm" 
  showDescription={false}
  className="custom-class"
/>`}</code></pre>
          
          <div class="mt-3 text-sm text-gray-600 dark:text-gray-400">
            <strong>Props:</strong>
            <ul class="list-disc list-inside mt-1 space-y-1">
              <li><code>statusCode</code>: number (required)</li>
              <li><code>size</code>: 'xs' | 'sm' | 'md' | 'lg'</li>
              <li><code>showDescription</code>: boolean</li>
              <li><code>className</code>: string</li>
            </ul>
          </div>
        </div>
        
        <div>
          <h3 class="font-medium text-gray-900 dark:text-white mb-3">HttpMethodBadge</h3>
          <pre class="bg-gray-100 dark:bg-gray-700 p-3 rounded text-xs overflow-x-auto"><code>{`<HttpMethodBadge 
  method="DELETE" 
  size="sm" 
  showDescription={false}
  className="custom-class"
/>`}</code></pre>
          
          <div class="mt-3 text-sm text-gray-600 dark:text-gray-400">
            <strong>Props:</strong>
            <ul class="list-disc list-inside mt-1 space-y-1">
              <li><code>method</code>: string (required)</li>
              <li><code>size</code>: 'xs' | 'sm' | 'md' | 'lg'</li>
              <li><code>showDescription</code>: boolean</li>
              <li><code>className</code>: string</li>
            </ul>
          </div>
        </div>
      </div>
      
      <div class="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-md">
        <h4 class="font-medium text-blue-900 dark:text-blue-300 mb-2">
          <i class="fas fa-lightbulb mr-2"></i>Benefits
        </h4>
        <ul class="text-sm text-blue-800 dark:text-blue-200 space-y-1">
          <li>• <strong>Consistent Colors:</strong> Same color scheme across all applications</li>
          <li>• <strong>Reusable:</strong> Use anywhere in the app without duplicating styles</li>
          <li>• <strong>Responsive:</strong> Multiple size options for different contexts</li>
          <li>• <strong>Accessible:</strong> Proper contrast ratios and semantic markup</li>
          <li>• <strong>Flexible:</strong> Support for custom methods and status codes</li>
        </ul>
      </div>
    </div>

    <!-- Endpoint List Components -->
    <div class="mt-8 col-span-1 lg:col-span-2">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          <i class="fas fa-list mr-2"></i>Endpoint List Components
        </h2>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          These components demonstrate how the badge components are integrated into larger UI elements for displaying API endpoints.
        </p>

        <!-- Individual Endpoint Items -->
        <div class="space-y-4">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">Individual Endpoint Items</h3>
          <div class="space-y-3">
            {#each mockEndpoints as endpoint}
              <EndpointListItem 
                {endpoint}
                on:edit={handleEndpointAction}
                on:delete={handleEndpointAction}
                on:toggle={handleEndpointAction}
                on:duplicate={handleEndpointAction}
                on:test={handleEndpointAction}
              />
            {/each}
          </div>
        </div>

        <!-- Full Endpoint List -->
        <div class="mt-8">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Complete Endpoint List</h3>
          <EndpointList 
            endpoints={mockEndpoints}
            on:edit={handleEndpointAction}
            on:delete={handleEndpointAction}
            on:toggle={handleEndpointAction}
            on:duplicate={handleEndpointAction}
            on:test={handleEndpointAction}
            on:create={handleEndpointAction}
            on:refresh={handleEndpointAction}
          />
        </div>

        <!-- Compact Version -->
        <div class="mt-8">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Compact Endpoint List</h3>
          <EndpointList 
            endpoints={mockEndpoints}
            compact={true}
            showMetrics={false}
            on:edit={handleEndpointAction}
            on:delete={handleEndpointAction}
            on:toggle={handleEndpointAction}
            on:duplicate={handleEndpointAction}
            on:test={handleEndpointAction}
            on:create={handleEndpointAction}
            on:refresh={handleEndpointAction}
          />
        </div>

        <!-- Usage Examples -->
        <div class="mt-6 bg-gray-50 dark:bg-gray-900/50 p-4 rounded-lg">
          <h4 class="font-medium text-gray-900 dark:text-white mb-3">Usage Examples</h4>
          
          <div class="space-y-4 text-sm">
            <div>
              <h5 class="font-medium text-gray-800 dark:text-gray-200 mb-2">EndpointListItem Component</h5>
              <pre class="bg-gray-800 text-gray-200 p-3 rounded overflow-x-auto"><code>{`<EndpointListItem 
  endpoint={{
    id: '1',
    method: 'GET',
    path: '/api/users',
    name: 'Get Users',
    status: 'active',
    statusCode: 200,
    responseTime: 45
  }}
  on:edit={handleEdit}
  on:delete={handleDelete}
  on:toggle={handleToggle}
/>`}</code></pre>
            </div>

            <div>
              <h5 class="font-medium text-gray-800 dark:text-gray-200 mb-2">EndpointList Component</h5>
              <pre class="bg-gray-800 text-gray-200 p-3 rounded overflow-x-auto"><code>{`<EndpointList 
  endpoints={endpointArray}
  showStatusCode={true}
  showMetrics={true}
  compact={false}
  on:edit={handleEdit}
  on:create={handleCreate}
/>`}</code></pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Back Navigation -->
    <div class="mt-8 text-center">
      <a 
        href="/demo" 
        class="inline-flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors"
      >
        <i class="fas fa-arrow-left mr-2"></i>
        Back to Demo Menu
      </a>
    </div>
  </div>
</div>
