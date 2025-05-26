<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let preRequestScript: string = '';
	export let testScript: string = '';
	export let activeScriptTab: string = 'pre-request';

	const dispatch = createEventDispatcher();

	function updateScript(type: 'pre-request' | 'test', content: string) {
		if (type === 'pre-request') {
			preRequestScript = content;
		} else {
			testScript = content;
		}
		dispatch('scriptChange', { preRequestScript, testScript });
	}

	function setActiveScriptTab(tab: string) {
		activeScriptTab = tab;
	}

	function formatScript() {
		// TODO: Implement script formatting
		console.log('Format script');
	}

	function clearScript() {
		if (activeScriptTab === 'pre-request') {
			preRequestScript = '';
		} else {
			testScript = '';
		}
		dispatch('scriptChange', { preRequestScript, testScript });
	}

	function insertSnippet(snippet: string) {
		const currentScript = activeScriptTab === 'pre-request' ? preRequestScript : testScript;
		const newScript = currentScript + (currentScript ? '\n\n' : '') + snippet;
		updateScript(activeScriptTab as 'pre-request' | 'test', newScript);
	}

	// Script snippets
	const preRequestSnippets = [
		{
			name: 'Set Environment Variable',
			code: `// Set environment variable
pm.environment.set('variableName', 'value');`
		},
		{
			name: 'Add Dynamic Header',
			code: `// Add dynamic header
pm.request.headers.add({
    key: 'X-Request-ID',
    value: Math.random().toString(36).substr(2, 9)
});`
		},
		{
			name: 'Set Timestamp',
			code: `// Set current timestamp
pm.environment.set('timestamp', Date.now());
pm.environment.set('iso_timestamp', new Date().toISOString());`
		},
		{
			name: 'Generate UUID',
			code: `// Generate UUID
function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}
pm.environment.set('uuid', generateUUID());`
		}
	];

	const testSnippets = [
		{
			name: 'Status Code Test',
			code: `// Test status code
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});`
		},
		{
			name: 'Response Time Test',
			code: `// Test response time
pm.test("Response time is less than 200ms", function () {
    pm.expect(pm.response.responseTime).to.be.below(200);
});`
		},
		{
			name: 'JSON Response Test',
			code: `// Test JSON response
pm.test("Response is JSON", function () {
    pm.response.to.be.json;
});

pm.test("Response has required fields", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('name');
});`
		},
		{
			name: 'Header Test',
			code: `// Test response headers
pm.test("Content-Type is application/json", function () {
    pm.expect(pm.response.headers.get("Content-Type")).to.include("application/json");
});`
		},
		{
			name: 'Save Response Data',
			code: `// Save response data to environment
if (pm.response.code === 200) {
    const responseJson = pm.response.json();
    pm.environment.set('user_id', responseJson.id);
    pm.environment.set('auth_token', responseJson.token);
}`
		}
	];

	$: currentSnippets = activeScriptTab === 'pre-request' ? preRequestSnippets : testSnippets;
	$: currentScript = activeScriptTab === 'pre-request' ? preRequestScript : testScript;
	$: placeholder = activeScriptTab === 'pre-request' 
		? `// Pre-request Script Example:
// Set environment variables
pm.environment.set('timestamp', Date.now());

// Add dynamic headers
pm.request.headers.add({
    key: 'X-Request-ID',
    value: Math.random().toString(36).substr(2, 9)
});

console.log('Pre-request script executed');`
		: `// Test Script Example:
// Test status code
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

// Test response time
pm.test("Response time is acceptable", function () {
    pm.expect(pm.response.responseTime).to.be.below(1000);
});

// Test JSON response
pm.test("Response has data", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.be.an('object');
});`;
</script>

<!-- Scripts section -->
<div role="tabpanel" aria-labelledby="scripts-tab" class="space-y-4">
	<div class="flex justify-between items-center mb-4">
		<h2 class="text-sm font-semibold theme-text-primary flex items-center">
			<i class="fas fa-code text-purple-500 mr-2"></i>
			Scripts
		</h2>
		<div class="flex items-center space-x-2">
			<button 
				class="text-sm text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200 flex items-center"
				title="View script documentation and examples"
				aria-label="View script documentation"
			>
				<i class="fas fa-question-circle text-xs mr-1"></i>
				Help
			</button>
		</div>
	</div>

	<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
		<!-- Script Type Navigation -->
		<div class="border-b theme-border">
			<nav class="flex space-x-0" role="tablist" aria-label="Script type navigation">
				<button 
					class="flex-1 py-3 px-4 text-sm font-medium transition-all duration-200 border-r theme-border {activeScriptTab === 'pre-request' 
						? 'bg-blue-600 text-white shadow-sm' 
						: 'theme-bg-secondary theme-text-secondary hover:bg-gray-100 dark:hover:bg-gray-600'}"
					title="Edit pre-request script that runs before sending the request"
					aria-label="Pre-request Script tab"
					role="tab"
					aria-selected={activeScriptTab === 'pre-request'}
					on:click={() => setActiveScriptTab('pre-request')}
				>
					<i class="fas fa-play-circle mr-2"></i>
					Pre-request Script
				</button>
				<button 
					class="flex-1 py-3 px-4 text-sm font-medium transition-all duration-200 {activeScriptTab === 'test' 
						? 'bg-blue-600 text-white shadow-sm' 
						: 'theme-bg-secondary theme-text-secondary hover:bg-gray-100 dark:hover:bg-gray-600'}"
					title="Edit test scripts that run after receiving the response"
					aria-label="Tests tab"
					role="tab"
					aria-selected={activeScriptTab === 'test'}
					on:click={() => setActiveScriptTab('test')}
				>
					<i class="fas fa-check-circle mr-2"></i>
					Tests
				</button>
			</nav>
		</div>

		<!-- Script Editor -->
		<div class="p-4 space-y-4">
			<div class="flex justify-between items-center">
				<label for="script-editor" class="flex items-center text-sm font-medium theme-text-secondary">
					<i class="fas fa-edit text-purple-400 mr-2"></i>
					JavaScript Code Editor
				</label>
				<div class="flex space-x-2">
					<button 
						class="text-xs bg-gray-100 dark:bg-gray-700 theme-text-secondary px-2 py-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
						title="Format and beautify JavaScript code"
						aria-label="Format code"
						on:click={formatScript}
					>
						<i class="fas fa-indent mr-1"></i>
						Format
					</button>
					<button 
						class="text-xs bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 px-2 py-1 rounded hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors duration-200"
						title="Clear editor content"
						aria-label="Clear editor"
						on:click={clearScript}
					>
						<i class="fas fa-eraser mr-1"></i>
						Clear
					</button>
				</div>
			</div>
			
			<div class="relative">
				<textarea
					id="script-editor"
					class="w-full h-64 theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary font-mono text-sm transition-all duration-200 resize-y"
					placeholder={placeholder}
					title="JavaScript code editor for {activeScriptTab === 'pre-request' ? 'pre-request' : 'test'} scripts"
					aria-label="Script editor textarea"
					value={currentScript}
					on:input={(e) => updateScript(activeScriptTab as 'pre-request' | 'test', e.currentTarget.value)}
				></textarea>
			</div>
		</div>
	</div>

	<!-- Code Snippets -->
	<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
		<div class="p-4 border-b theme-border">
			<h3 class="text-sm font-medium theme-text-primary flex items-center">
				<i class="fas fa-code text-green-400 mr-2"></i>
				{activeScriptTab === 'pre-request' ? 'Pre-request' : 'Test'} Snippets
			</h3>
			<p class="text-xs theme-text-muted mt-1">
				Click on any snippet to add it to your script
			</p>
		</div>
		<div class="p-4">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
				{#each currentSnippets as snippet}
					<button
						class="text-left p-3 bg-gray-50 dark:bg-gray-900 rounded-md border theme-border hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors duration-200"
						title="Add {snippet.name} to script"
						aria-label="Insert {snippet.name} snippet"
						on:click={() => insertSnippet(snippet.code)}
					>
						<div class="flex items-center justify-between mb-1">
							<span class="text-sm font-medium theme-text-primary">{snippet.name}</span>
							<i class="fas fa-plus text-xs theme-text-muted"></i>
						</div>
						<code class="text-xs theme-text-muted block truncate">
							{snippet.code.split('\n')[1]?.trim() || snippet.code.split('\n')[0]}
						</code>
					</button>
				{/each}
			</div>
		</div>
	</div>

	<!-- Script API Reference -->
	<div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
		<h4 class="text-sm font-medium theme-text-primary mb-2 flex items-center">
			<i class="fas fa-book text-blue-400 mr-2"></i>
			Script API Reference
		</h4>
		<div class="text-xs theme-text-secondary space-y-1">
			<div><code class="bg-white dark:bg-gray-800 px-1 rounded">pm.environment.set(key, value)</code> - Set environment variable</div>
			<div><code class="bg-white dark:bg-gray-800 px-1 rounded">pm.environment.get(key)</code> - Get environment variable</div>
			<div><code class="bg-white dark:bg-gray-800 px-1 rounded">pm.request.headers.add(header)</code> - Add request header</div>
			<div><code class="bg-white dark:bg-gray-800 px-1 rounded">pm.response.status</code> - Response status code</div>
			<div><code class="bg-white dark:bg-gray-800 px-1 rounded">pm.response.json()</code> - Parse response as JSON</div>
			<div><code class="bg-white dark:bg-gray-800 px-1 rounded">pm.test(name, function)</code> - Create a test</div>
		</div>
	</div>
</div>
