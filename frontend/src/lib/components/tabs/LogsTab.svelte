<script lang="ts">
	import type { Project } from "$lib/api/mockoonApi";

  interface Log {
    method: string;
    path: string;
    timestamp: string;
    request: {
      headers: Record<string, string>;
      query?: Record<string, string>;
      body?: any;
    };
    response: {
      status: number;
      headers: Record<string, string>;
      body: any;
    };
  }

  interface Config {
    uuid: string;
    name: string;
    configFile: string;
    port: number;
    url: string;
    size: string;
    modified: string;
    inUse: boolean;
  }

  export let selectedProject: Project;

  // Dummy data for logs
  let logs: Log[] = [
    {
      method: 'GET',
      path: '/api/v1/users',
      timestamp: '2024-03-20 10:00:00',
      request: {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer token123'
        },
        query: {
          page: '1',
          limit: '10'
        }
      },
      response: {
        status: 200,
        headers: {
          'Content-Type': 'application/json',
          'X-Total-Count': '2'
        },
        body: {
          users: [
            { id: 1, name: 'John Doe', email: 'john@example.com' },
            { id: 2, name: 'Jane Smith', email: 'jane@example.com' }
          ],
          total: 2
        }
      }
    },
    {
      method: 'POST',
      path: '/api/v1/users',
      timestamp: '2024-03-20 10:01:00',
      request: {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer token123'
        },
        body: {
          name: 'New User',
          email: 'new@example.com'
        }
      },
      response: {
        status: 201,
        headers: {
          'Content-Type': 'application/json',
          'Location': '/api/v1/users/3'
        },
        body: {
          id: 3,
          name: 'New User',
          email: 'new@example.com',
          createdAt: '2024-03-20T10:01:00Z'
        }
      }
    }
  ];
</script>

<div class="w-full bg-gray-800 p-4">
  <div class="bg-gray-700 p-4 rounded mb-4 flex items-center">
    <i class="fas fa-info-circle text-blue-500 text-2xl mr-2"></i>
    <span class="text-xl font-bold text-blue-500">
      Logs for: {selectedProject.name}
    </span>
  </div>
  <div class="flex items-center bg-gray-700 p-2 rounded mb-4">
    <i class="fas fa-search text-white text-lg mr-2"></i>
    <input
      type="text"
      id="log-search"
      placeholder="Search Logs"
      class="w-full bg-gray-700 text-white py-1 px-2 rounded text-sm"
    />
  </div>
  <div class="space-y-4">
    {#each logs as log}
      <div class="bg-gray-700 p-4 rounded">
        <div class="flex justify-between items-center mb-2">
          <span class="text-sm font-bold">
            <strong>{log.method}</strong> {log.path}
          </span>
          <span class="text-xs text-gray-400">{log.timestamp}</span>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <h3 class="text-sm font-semibold mb-2">Request</h3>
            <pre class="bg-gray-800 p-2 rounded text-xs">{JSON.stringify(log.request, null, 2)}</pre>
          </div>
          <div>
            <h3 class="text-sm font-semibold mb-2">Response</h3>
            <pre class="bg-gray-800 p-2 rounded text-xs">{JSON.stringify(log.response, null, 2)}</pre>
          </div>
        </div>
      </div>
    {/each}
  </div>
</div> 