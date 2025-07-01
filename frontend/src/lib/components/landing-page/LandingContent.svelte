<script lang="ts">
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { toast } from '$lib/stores/toast';
	import RecentProjects from '$lib/components/common/RecentProjects.svelte';
	import ProjectSearchResults from '$lib/components/common/ProjectSearchResults.svelte';
	import { publicConfig } from '$lib/stores/publicConfig';
	import { onMount } from 'svelte';
	import { checkAliasAndSearchProjects, type ProjectSearchResult } from '$lib/api/BeoApi';
	
	let projectName = '';
	let isLoading = false;
	let searchResults: ProjectSearchResult[] = [];
	let showSearchResults = false;
	let searchTimeout: ReturnType<typeof setTimeout>;
	let isSearching = false;
	let aliasAvailable = true; // Track if alias is available

	// Computed property for URL format display
	$: urlFormatDisplay = getUrlFormatDisplay(projectName, $publicConfig?.mock_url_format || 'subdomain');

	// Features data for the landing page
	const features = [
		{
			title: 'Mock Server',
			description:
				'Create custom API endpoints with configurable responses. Perfect for frontend development and testing.',
			icon: 'fas fa-server',
			color: 'text-green-400',
			bgColor: 'bg-green-400/20'
		},
		{
			title: 'Proxy Mode',
			description:
				'Smart proxy that uses mocks when available, otherwise forwards requests to real endpoints.',
			icon: 'fas fa-exchange-alt',
			color: 'text-blue-400',
			bgColor: 'bg-blue-400/20'
		},
		{
			title: 'Request Forwarder',
			description:
				'Always forward requests to target endpoints while logging all traffic for analysis.',
			icon: 'fas fa-arrow-right',
			color: 'text-purple-400',
			bgColor: 'bg-purple-400/20'
		},
		{
			title: 'Request Logging',
			description:
				'Comprehensive logging of all requests and responses with filtering and search capabilities.',
			icon: 'fas fa-list-alt',
			color: 'text-yellow-400',
			bgColor: 'bg-yellow-400/20'
		},
		{
			title: 'Response Templates',
			description: 'Create dynamic responses with templates, rules, and conditional logic.',
			icon: 'fas fa-code',
			color: 'text-orange-400',
			bgColor: 'bg-orange-400/20'
		},
		{
			title: 'Multi-User Workspaces',
			description: 'Collaborate with your team using shared workspaces and project management.',
			icon: 'fas fa-users',
			color: 'text-indigo-400',
			bgColor: 'bg-indigo-400/20'
		}
	];

	const modes = [
		{
			name: 'Mock',
			description: 'Serves predefined mock responses only',
			icon: 'fas fa-server',
			color: 'bg-green-600'
		},
		{
			name: 'Proxy',
			description: 'Uses mocks when available, otherwise forwards requests',
			icon: 'fas fa-exchange-alt',
			color: 'bg-blue-600'
		},
		{
			name: 'Forwarder',
			description: 'Always forwards all requests to target endpoint',
			icon: 'fas fa-arrow-right',
			color: 'bg-purple-600'
		},
		{
			name: 'Disabled',
			description: 'Endpoint inactive - no responses served',
			icon: 'fas fa-ban',
			color: 'bg-gray-600'
		}
	];

	// Authentication state
	$: authenticated = $isAuthenticated;

	// Function to generate URL format display based on configuration
	function getUrlFormatDisplay(alias: string, format: string): string {
		const cleanAlias = alias.trim() || 'alias';
		
		// Backend sends the exact format, e.g.:
		// - "alias.localhost:3600" for subdomain mode
		// - "localhost:3600/alias" for path mode
		if (format && format.includes('alias')) {
			return format.replace('alias', cleanAlias);
		}
		
		// Fallback if format is just "subdomain" or "path" (for backward compatibility)
		if (format === 'subdomain') {
			return `${cleanAlias}.localhost:3600`;
		} else {
			return `localhost:3600/${cleanAlias}`;
		}
	}



	async function createProject() {
		if (!projectName.trim()) {
			toast.error('Please enter a project alias');
			return;
		}

		isLoading = true;
		try {
			// For now, redirect to home to create project
			// In the future, we can implement inline project creation
			await goto('/home');
		} catch (err) {
			toast.error(err);
		} finally {
			isLoading = false;
		}
	}

	async function handleLogin() {
		await goto('/login');
	}

	// Debounced search function
	function handleSearchInput() {
		// Clear previous timeout
		if (searchTimeout) {
			clearTimeout(searchTimeout);
		}

		// Hide results if input is empty
		if (!projectName.trim()) {
			showSearchResults = false;
			searchResults = [];
			aliasAvailable = true; // Reset to available when empty
			return;
		}

		// Clear search results while searching but don't reset aliasAvailable yet
		showSearchResults = false;
		searchResults = [];

		// Set timeout for 500ms
		searchTimeout = setTimeout(async () => {
			await searchProjects(projectName);
		}, 500);
	}

	// Search for existing projects
	async function searchProjects(query: string) {
		if (!query.trim() || !authenticated) return;

		isSearching = true;
		try {
			const response = await checkAliasAndSearchProjects(query.trim());
			if (response.projects) {
				searchResults = response.projects;
				showSearchResults = searchResults.length > 0;
				aliasAvailable = response.available; // Update availability state
			} else {
				searchResults = [];
				showSearchResults = false;
				aliasAvailable = true; // Default to available if no response
			}
		} catch (error) {
			console.error('Search error:', error);
			toast.error('Failed to search projects');
			searchResults = [];
			showSearchResults = false;
			aliasAvailable = true; // Default to available on error
		} finally {
			isSearching = false;
		}
	}

	// Handle project selection from search results
	function handleProjectSelect(project: ProjectSearchResult) {
		// Navigate to project management page like recent projects
		goto(`/home/workspace/${project.workspace_id}/projects/${project.id}`);
		
		// Clear search
		projectName = '';
		showSearchResults = false;
		searchResults = [];
		aliasAvailable = true; // Reset availability
	}

	// Close search results when clicking outside
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as Element;
		if (target && !target.closest('.search-container')) {
			showSearchResults = false;
		}
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
			if (searchTimeout) {
				clearTimeout(searchTimeout);
			}
		};
	});
</script>

<!-- Main Content -->
<main class="flex-1">

	<!-- Hero Section -->
	<section class="bg-gradient-to-b from-blue-50 to-white dark:from-gray-800 dark:to-gray-900 py-12 pt-24">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center">
				<h1 class="text-3xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4">
					Need instant API endpoints?
				</h1>
				<h2 class="text-xl md:text-2xl font-light text-gray-700 dark:text-gray-300 mb-6">
					Create a <span class="text-indigo-600 dark:text-indigo-400 font-semibold"
						>mock server in seconds</span
					>
				</h2>
				<p class="text-base text-gray-600 dark:text-gray-400 mb-6 max-w-2xl mx-auto">
					<i class="fas fa-lightning-bolt text-indigo-600 dark:text-indigo-400 mr-2"></i>
					Zero setup, instant deployment, maximum productivity.
				</p>

				<!-- Quick Action Buttons -->
				<div class="flex flex-col sm:flex-row gap-3 justify-center items-center mb-8">
					<button
						on:click={() => {
							document.getElementById('quick-deploy')?.scrollIntoView({ behavior: 'smooth' });
						}}
						class="bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 dark:from-blue-600 dark:to-purple-700 dark:hover:from-blue-700 dark:hover:to-purple-800 text-white py-3 px-6 rounded-lg text-base font-medium transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-1 flex items-center"
						title="Deploy instantly with Docker - One command to run!"
						aria-label="Deploy instantly with Docker"
					>
						<i class="fab fa-docker mr-2 text-lg"></i>
						🚀  Deploy in Seconds
					</button>
					
					{#if !authenticated}
						<button
							on:click={handleLogin}
							class="bg-white dark:bg-gray-800 border-2 border-gray-200 dark:border-gray-600 hover:border-blue-500 dark:hover:border-blue-400 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 py-3 px-6 rounded-lg text-base font-medium transition-all duration-300 flex items-center shadow-md hover:shadow-lg"
							title="Login to create cloud projects"
							aria-label="Login to create cloud projects"
						>
							<i class="fas fa-cloud mr-2"></i>
							Try Cloud Version
						</button>
					{/if}
				</div>

				<!-- Main Dashboard Section - Side by side layout -->
				<div class="max-w-6xl mx-auto">
					<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
						<!-- Create New Mock Server -->
						{#if authenticated}
							<div class="bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-xl border border-gray-200 dark:border-gray-700 p-6">
								<h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
									<i class="fas fa-rocket text-indigo-600 dark:text-indigo-400 mr-2"></i>
									Create Your Next Mock Server
								</h3>

								<div class="search-container relative">
									<div class="flex flex-col sm:flex-row gap-3 mb-4">
										<!-- Use a single input with placeholder based on format -->
										<div class="flex-1 flex focus-within:ring-2 focus-within:ring-indigo-500 rounded-lg">
											<input
												bind:value={projectName}
												on:input={handleSearchInput}
												type="text"
												placeholder="your-project-alias"
												class="flex-1 px-3 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none transition-colors text-sm"
												title="Enter a alias for your mock server project"
												aria-label="Project alias input"
											/>
										</div>
									</div>

									<!-- Search Results Component -->
									<ProjectSearchResults 
										{searchResults}
										showResults={showSearchResults}
										onProjectSelect={handleProjectSelect}
									/>

									<!-- Search Loading Indicator -->
									{#if isSearching}
										<div class="absolute right-3 top-[58px] text-gray-400">
											<i class="fas fa-spinner fa-spin text-sm"></i>
										</div>
									{/if}
								</div>

								<p class="text-xs text-gray-600 dark:text-gray-400 mb-4">
									{#if showSearchResults && searchResults.length > 0}
										<span class="text-orange-600 dark:text-orange-400">
											<i class="fas fa-info-circle mr-1"></i>
											Found existing projects. Click to open or continue typing to create new.
										</span>
									{:else if projectName.trim() && !aliasAvailable}
										<span class="text-red-600 dark:text-red-400">
											<i class="fas fa-times-circle mr-1"></i>
											Alias "{projectName}" is already used by another project
										</span>
									{:else if $publicConfig?.mock_url_format === 'subdomain'}
										Your mock server will be available at: <span class="font-mono text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20 px-2 py-1 rounded">{urlFormatDisplay}</span>
									{:else}
										Your mock server will be available at: <span class="font-mono text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20 px-2 py-1 rounded">{urlFormatDisplay}</span>
									{/if}
								</p>

								<button
									on:click={createProject}
									disabled={isLoading || (!aliasAvailable && !!projectName.trim())}
									class="w-full bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed text-white py-2.5 px-4 rounded-lg font-medium transition-all duration-200 flex items-center justify-center shadow-lg hover:shadow-xl text-sm"
									title={(!aliasAvailable && !!projectName.trim()) ? "Alias is not available" : "Create new mock server project"}
									aria-label={(!aliasAvailable && !!projectName.trim()) ? "Alias is not available" : "Create new mock server project"}
								>
									{#if isLoading}
										<i class="fas fa-spinner fa-spin mr-2"></i>
										Setting up your server...
									{:else if !aliasAvailable && projectName.trim()}
										<i class="fas fa-exclamation-triangle mr-2"></i>
										Alias Not Available
									{:else}
										<i class="fas fa-plus-circle mr-2"></i>
										Create Mock Server
									{/if}
								</button>
							</div>
						{:else}
							<div class="bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-xl border border-gray-200 dark:border-gray-700 p-6 text-center">
								<h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
									<i class="fas fa-user-circle text-indigo-600 dark:text-indigo-400 mr-2"></i>
									Get Started
								</h3>
								<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
									Login to create and manage your mock servers
								</p>
								<button
									on:click={handleLogin}
									class="w-full bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white py-2.5 px-4 rounded-lg font-medium transition-all duration-200 flex items-center justify-center shadow-lg hover:shadow-xl text-sm"
									title="Login to create cloud projects"
									aria-label="Login to create cloud projects"
								>
									<i class="fas fa-sign-in-alt mr-2"></i>
									Login to Continue
								</button>
								<p class="text-xs text-gray-500 dark:text-gray-400 mt-3">
									Free to use • No credit card required
								</p>
							</div>
						{/if}

						<!-- Recent Projects Section -->
						<div class="bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
							<div class="p-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-750/50">
								<h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
									<i class="fas fa-history text-indigo-600 dark:text-indigo-400 mr-2"></i>
									Recent Projects
								</h3>
								<p class="text-xs text-gray-600 dark:text-gray-400 mt-1">
									Quick access to your recently used mock servers
								</p>
							</div>
							
							<div class="p-1">
								<RecentProjects 
									showTitle={false} 
									maxItems={4}
									onProjectSelect={(project) => {
										// Handle project selection
										if (authenticated) {
											goto(`/home/workspace/${project.workspaceId}/projects/${project.id}`);
										} else {
											toast.info('Please login to access your projects');
										}
									}}
								/>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Beo Echo Modes Section -->
	<section id="modes" class="py-16 bg-gray-50 dark:bg-gray-800">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center mb-12">
				<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">
					<i class="fas fa-cogs text-blue-600 mr-2"></i>
					Beo Echo Operating Modes
				</h2>
				<p class="text-lg text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
					Choose the perfect mode for your development workflow
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
				{#each modes as mode}
					<div
						class="bg-white dark:bg-gray-700 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
					>
						<div class="flex items-center mb-4">
							<div class="w-10 h-10 {mode.color} rounded-lg flex items-center justify-center mr-3">
								<i class="{mode.icon} text-white"></i>
							</div>
							<h3 class="font-semibold text-gray-900 dark:text-white">{mode.name}</h3>
						</div>
						<p class="text-gray-600 dark:text-gray-300 text-sm">
							{mode.description}
						</p>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section id="features" class="py-16">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center mb-12">
				<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">
					<i class="fas fa-star text-blue-600 mr-2"></i>
					Beo Echo Features & Use Cases
				</h2>
				<p class="text-lg text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">
					Discover how Beo Echo can streamline your development workflow, speed up API integrations
					and software delivery.
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
				{#each features as feature}
					<div
						class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
					>
						<div class="flex items-center mb-4">
							<div
								class="w-12 h-12 {feature.bgColor} rounded-lg flex items-center justify-center mr-4"
							>
								<i class="{feature.icon} {feature.color} text-xl"></i>
							</div>
							<h3 class="text-lg font-semibold text-gray-900 dark:text-white">{feature.title}</h3>
						</div>
						<p class="text-gray-600 dark:text-gray-300">
							{feature.description}
						</p>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Use Cases Section -->
	<section class="py-16 bg-gray-50 dark:bg-gray-800">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<!-- Replace Dependencies -->
				<div class="bg-white dark:bg-gray-700 rounded-lg shadow-md p-8">
					<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
						<i class="fas fa-server text-green-600 mr-2"></i>
						Replace Dependencies in Tests
					</h3>
					<p class="text-gray-600 dark:text-gray-300 mb-6">
						Mock external APIs and services in your tests for faster, more reliable test suites.
					</p>
					<button
						class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
						title="Create a mock API server"
						aria-label="Create a mock API server"
					>
						<i class="fas fa-plus mr-2"></i>
						Create Mock Server
					</button>
				</div>

				<!-- Start Integration -->
				<div class="bg-white dark:bg-gray-700 rounded-lg shadow-md p-8">
					<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
						<i class="fas fa-code text-blue-600 mr-2"></i>
						Start Integration Before APIs are Ready
					</h3>
					<p class="text-gray-600 dark:text-gray-300 mb-6">
						Begin frontend development immediately with mock APIs that match your planned backend
						structure.
					</p>
					<button
						class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
						title="Create a mock API server"
						aria-label="Create a mock API server"
					>
						<i class="fas fa-plus mr-2"></i>
						Create Mock Server
					</button>
				</div>

				<!-- Public HTTP Endpoint -->
				<div class="bg-white dark:bg-gray-700 rounded-lg shadow-md p-8">
					<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
						<i class="fas fa-globe text-purple-600 mr-2"></i>
						Create Public HTTP Endpoints
					</h3>
					<p class="text-gray-600 dark:text-gray-300 mb-6">
						Generate publicly accessible endpoints for webhooks, API testing, or sharing with
						teammates.
					</p>
					<button
						class="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
						title="View HTTP request details"
						aria-label="View HTTP request details"
					>
						<i class="fas fa-eye mr-2"></i>
						View Requests
					</button>
				</div>

				<!-- Partial Mocks -->
				<div class="bg-white dark:bg-gray-700 rounded-lg shadow-md p-8">
					<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
						<i class="fas fa-puzzle-piece text-green-600 mr-2"></i>
						Partial Mocks
					</h3>
					<p class="text-gray-600 dark:text-gray-300 mb-6">
						Mock only specific endpoints while forwarding others to real services for hybrid testing
						scenarios.
					</p>
					<button
						class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
						title="Create a partial mock server"
						aria-label="Create a partial mock server"
					>
						<i class="fas fa-plus mr-2"></i>
						Create Mock Server
					</button>
				</div>
			</div>
		</div>
	</section>

	<!-- Quick Deploy Section -->
	<section id="quick-deploy" class="py-16 bg-gradient-to-r from-blue-500 to-indigo-600 dark:from-blue-700 dark:to-purple-800">
		<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
			<div class="mb-8">
				<h2 class="text-3xl font-bold text-white mb-4">
					<i class="fab fa-docker text-white mr-3"></i>
					Deploy in Seconds
				</h2>
				<p class="text-xl text-blue-50 dark:text-blue-100 mb-6">
					One command to run Beo Echo locally with Docker
				</p>
			</div>

			<!-- Docker Command -->
			<div class="bg-gray-800 dark:bg-gray-900 rounded-lg p-6 mb-8 text-left overflow-x-auto border border-gray-700 dark:border-gray-600">
				<div class="flex items-center justify-between mb-3">
					<span class="text-green-400 dark:text-green-300 text-sm font-mono">Terminal</span>
					<button
						class="bg-blue-500 hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700 text-white px-3 py-1 rounded text-xs font-medium transition-colors"
						title="Copy Docker command to clipboard"
						aria-label="Copy Docker command to clipboard"
						on:click={() => {
							navigator.clipboard.writeText('docker run -d --platform linux/amd64 -p 8080:80 -v $(pwd)/beo-echo-config:/app/configs/ ghcr.io/yogasw/beo-echo:latest');
							toast.success('Docker command copied to clipboard!');
						}}
					>
						<i class="fas fa-copy mr-1"></i>
						Copy
					</button>
				</div>
				<code class="text-green-400 dark:text-green-300 font-mono text-sm block leading-relaxed">
					<span class="text-gray-400 dark:text-gray-500">$</span> docker run -d --platform linux/amd64 -p 8080:80 \<br>
					<span class="ml-4">-v $(pwd)/beo-echo-config:/app/configs/ \</span><br>
					<span class="ml-4">ghcr.io/yogasw/beo-echo:latest</span>
				</code>
			</div>

			<!-- Quick Steps -->
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
				<div class="bg-white/15 dark:bg-white/10 backdrop-blur-sm rounded-lg p-6 border border-white/20 dark:border-white/10">
					<div class="w-12 h-12 bg-white/25 dark:bg-white/20 rounded-full flex items-center justify-center mx-auto mb-4">
						<span class="text-2xl font-bold text-white">1</span>
					</div>
					<h3 class="text-lg font-semibold text-white mb-2">Run Command</h3>
					<p class="text-blue-50 dark:text-blue-100 text-sm">
						Execute the Docker command in your terminal
					</p>
				</div>

				<div class="bg-white/15 dark:bg-white/10 backdrop-blur-sm rounded-lg p-6 border border-white/20 dark:border-white/10">
					<div class="w-12 h-12 bg-white/25 dark:bg-white/20 rounded-full flex items-center justify-center mx-auto mb-4">
						<span class="text-2xl font-bold text-white">2</span>
					</div>
					<h3 class="text-lg font-semibold text-white mb-2">Open Browser</h3>
					<p class="text-blue-50 dark:text-blue-100 text-sm">
						Access at <span class="font-mono bg-white/25 dark:bg-white/20 px-1 rounded">localhost:8080</span>
					</p>
				</div>

				<div class="bg-white/15 dark:bg-white/10 backdrop-blur-sm rounded-lg p-6 border border-white/20 dark:border-white/10">
					<div class="w-12 h-12 bg-white/25 dark:bg-white/20 rounded-full flex items-center justify-center mx-auto mb-4">
						<span class="text-2xl font-bold text-white">3</span>
					</div>
					<h3 class="text-lg font-semibold text-white mb-2">Login & Start</h3>
					<p class="text-blue-50 dark:text-blue-100 text-sm">
						Use <span class="font-mono bg-white/25 dark:bg-white/20 px-1 rounded">admin@admin.com</span> / <span class="font-mono bg-white/25 dark:bg-white/20 px-1 rounded">admin</span>
					</p>
				</div>
			</div>

			<!-- Additional Info -->
			<div class="bg-white/15 dark:bg-white/10 backdrop-blur-sm rounded-lg p-6 text-left border border-white/20 dark:border-white/10">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<div>
						<h4 class="text-white font-semibold mb-3 flex items-center">
							<i class="fas fa-database text-blue-100 dark:text-blue-200 mr-2"></i>
							Database Options
						</h4>
						<ul class="text-blue-50 dark:text-blue-100 text-sm space-y-1">
							<li>• Default: SQLite (auto-created)</li>
							<li>• PostgreSQL: Set <span class="font-mono bg-white/25 dark:bg-white/20 px-1 rounded">DATABASE_URL</span></li>
						</ul>
					</div>
					<div>
						<h4 class="text-white font-semibold mb-3 flex items-center">
							<i class="fas fa-cog text-blue-100 dark:text-blue-200 mr-2"></i>
							Configuration
						</h4>
						<ul class="text-blue-50 dark:text-blue-100 text-sm space-y-1">
							<li>• Config stored in <span class="font-mono bg-white/25 dark:bg-white/20 px-1 rounded">./beo-echo-config/</span></li>
							<li>• Persistent data across container restarts</li>
						</ul>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Pricing Section -->
	<section
		id="pricing"
		class="py-16 bg-gradient-to-b from-gray-50 to-white dark:from-gray-800 dark:to-gray-900"
	>
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center mb-12">
				<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">
					<i class="fas fa-tags text-blue-600 mr-2"></i>
					Simple, Transparent Pricing
				</h2>
				<p class="text-lg text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
					Choose the plan that fits your needs. Start with Community Edition for free, or unlock
					advanced features with Cloud.
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-7xl mx-auto">
				<!-- Community Edition -->
				<div
					class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg border border-gray-200 dark:border-gray-700 p-8 relative"
				>
					<div class="mb-8">
						<h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Community Edition</h3>
						<p class="text-gray-600 dark:text-gray-400 mb-4">
							Perfect for individual developers and small projects
						</p>
						<div class="flex items-baseline">
							<span class="text-3xl font-bold text-gray-900 dark:text-white">Free</span>
							<span class="text-gray-500 dark:text-gray-400 ml-2">forever</span>
						</div>
					</div>

					<div class="space-y-4 mb-8">
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Unlimited mock servers</span>
						</div>
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Request logging & filtering</span>
						</div>
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Multi-user workspaces</span>
						</div>
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Docker deployment</span>
						</div>
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">SQLite & PostgreSQL support</span>
						</div>
					</div>

					<button
						class="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 px-6 rounded-lg font-medium transition-colors flex items-center justify-center"
						title="Deploy Beo Echo with Docker"
						aria-label="Deploy Beo Echo with Docker"
						on:click={() => {
							// Scroll to deploy section
							document.getElementById('quick-deploy')?.scrollIntoView({ behavior: 'smooth' });
						}}
					>
						<i class="fab fa-docker mr-2"></i>
						Deploy Now
					</button>
				</div>

				<!-- Cloud Plan -->
				<div
					class="bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-2xl shadow-lg border-2 border-blue-200 dark:border-blue-700 p-8 relative"
				>
					<!-- Popular Badge -->
					<div class="absolute -top-4 left-1/2 transform -translate-x-1/2">
						<span class="bg-blue-600 text-white px-4 py-1 rounded-full text-sm font-medium">
							Most Popular
						</span>
					</div>

					<div class="mb-8">
						<h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Cloud</h3>
						<p class="text-gray-600 dark:text-gray-400 mb-4">
							For teams and production deployments
						</p>
						<div class="flex items-baseline">
							<span class="text-3xl font-bold text-gray-900 dark:text-white">$0</span>
							<span class="text-gray-500 dark:text-gray-400 ml-2">per month</span>
						</div>
					</div>

					<div class="space-y-4 mb-8">
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Everything in Community</span>
						</div>
					</div>

					<button
						class="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 px-6 rounded-lg font-medium transition-colors"
						title="Get started with Cloud edition"
						aria-label="Get started with Cloud edition"
					>
						<i class="fas fa-rocket mr-2"></i>
						Get Started
					</button>
				</div>

				<!-- Pro Plan -->
				<div
					class="bg-gray-50 dark:bg-gray-800 rounded-2xl shadow-lg border border-gray-200 dark:border-gray-700 p-8 relative opacity-75"
				>
					<!-- Coming Soon Badge -->
					<div class="absolute -top-4 left-1/2 transform -translate-x-1/2">
						<span class="bg-gray-600 text-white px-4 py-1 rounded-full text-sm font-medium">
							Coming Soon
						</span>
					</div>

					<div class="mb-8">
						<h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Pro</h3>
						<p class="text-gray-600 dark:text-gray-400 mb-4">
							For enterprise and advanced use cases
						</p>
						<div class="flex items-baseline">
							<span class="text-3xl font-bold text-gray-900 dark:text-white">Contact Us</span>
							<span class="text-gray-500 dark:text-gray-400 ml-2">for pricing</span>
						</div>
					</div>

					<div class="space-y-4 mb-8">
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Everything in Cloud</span>
						</div>
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Advanced analytics</span>
						</div>

						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Custom domains</span>
						</div>

						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Priority support</span>
						</div>
					</div>

					<button
						disabled
						class="w-full bg-gray-400 text-gray-600 py-3 px-6 rounded-lg font-medium cursor-not-allowed"
						title="Pro plan coming soon"
						aria-label="Pro plan coming soon"
					>
						Coming Soon
					</button>
				</div>
			</div>
		</div>
	</section>
</main>
