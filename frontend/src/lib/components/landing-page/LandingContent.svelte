<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { getProjects } from '$lib/api/BeoApi';
	import { toast } from '$lib/stores/toast';
	import type { Project } from '$lib/api/BeoApi';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';

	let recentProjects: Project[] = [];
	let projectName = '';
	let isLoading = false;
	let loadingRecentProjects = false;

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

	// Load recent projects if authenticated
	export let authenticated = false;
	$: if (authenticated) {
		loadRecentProjects();
	}

	async function loadRecentProjects() {
		if (!authenticated) return;

		try {
			loadingRecentProjects = true;
			const projectsData = await getProjects();
			// Show only the 3 most recent projects
			recentProjects = projectsData.slice(0, 3);
		} catch (err) {
			console.error('Failed to load recent projects:', err);
		} finally {
			loadingRecentProjects = false;
		}
	}

	async function createProject() {
		if (!projectName.trim()) {
			toast.error('Please enter a project name');
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

	function generateProjectUrl(alias: string) {
		return `https://${alias}.beo-echo.dev`;
	}
</script>

<!-- Main Content -->
<main class="flex-1">
	<!-- Hero Section -->
	<section class="bg-gradient-to-b from-blue-50 to-white dark:from-gray-800 dark:to-gray-900 py-20">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center">
				<h1 class="text-4xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6">
					Unfinished APIs slowing you down?
				</h1>
				<h2 class="text-2xl md:text-4xl font-light text-gray-700 dark:text-gray-300 mb-8">
					Deploy a <span class="text-blue-600 dark:text-blue-400 font-semibold"
						>mock API in a few seconds</span
					>
				</h2>
				<p class="text-lg text-gray-600 dark:text-gray-400 mb-8 max-w-2xl mx-auto">
					<i class="fas fa-exchange-alt text-blue-600 mr-2"></i>
					No downloads, No dependencies, No Delays.
				</p>

				<!-- Quick Start Section -->
				{#if authenticated}
					<div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-8 max-w-2xl mx-auto mb-8">
						<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-6">
							<i class="fas fa-rocket text-blue-600 mr-2"></i>
							Launch a mock server now!
						</h3>

						<div class="flex flex-col sm:flex-row gap-4 mb-6">
							<div class="flex-1">
								<input
									bind:value={projectName}
									type="text"
									placeholder="Project Name"
									class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
								/>
							</div>
							<div class="flex items-center text-gray-500 dark:text-gray-400 text-sm">
								.beo-echo.dev
							</div>
						</div>

						<p class="text-sm text-gray-600 dark:text-gray-400 mb-6">
							A sub-domain will be created where you can send HTTP or API requests.
						</p>

						<button
							on:click={createProject}
							disabled={isLoading}
							class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white py-3 px-6 rounded-lg font-medium transition-colors flex items-center justify-center"
							title="Create new mock server"
							aria-label="Create new mock server"
						>
							{#if isLoading}
								<i class="fas fa-spinner fa-spin mr-2"></i>
								Creating...
							{:else}
								<i class="fas fa-plus mr-2"></i>
								Create Mock Server
							{/if}
						</button>

						{#if recentProjects.length > 0}
							<div class="mt-8">
								<h4 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-4">
									Your recent endpoints:
								</h4>
								<div class="space-y-2">
									{#each recentProjects as project}
										<div
											class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg"
										>
											<div class="flex items-center">
												<div
													class="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded flex items-center justify-center mr-3"
												>
													<i class="fas fa-server text-blue-600 dark:text-blue-400 text-sm"></i>
												</div>
												<div>
													<p class="font-medium text-gray-900 dark:text-white text-sm">
														{project.name}
													</p>
													<p class="text-xs text-gray-500 dark:text-gray-400">
														{generateProjectUrl(project.alias)}
													</p>
												</div>
											</div>
											<button
												class="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 text-sm font-medium"
												title="Open project"
												aria-label="Open project {project.name}"
											>
												Open
											</button>
										</div>
									{/each}
								</div>
							</div>
						{:else if loadingRecentProjects}
							<div class="mt-8">
								<SkeletonLoader type="list" count={3} />
							</div>
						{/if}
					</div>
				{:else}
					<div class="space-y-4">
						<button
							on:click={handleLogin}
							class="bg-blue-600 hover:bg-blue-700 text-white py-4 px-8 rounded-lg text-lg font-medium transition-colors"
							title="Login to get started"
							aria-label="Login to get started"
						>
							<i class="fas fa-sign-in-alt mr-2"></i>
							Get Started - Login
						</button>
						<p class="text-sm text-gray-600 dark:text-gray-400">
							Free to use • No credit card required
						</p>
					</div>
				{/if}

				<!-- Demo Video Button -->
				<div class="mt-8">
					<button
						class="inline-flex items-center text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 font-medium"
						title="Watch quick demo"
						aria-label="Watch quick demo"
					>
						<div
							class="flex items-center justify-center w-12 h-12 bg-blue-100 dark:bg-blue-900 rounded-full mr-3"
						>
							<i class="fas fa-play text-blue-600 dark:text-blue-400"></i>
						</div>
						Check out a quick demo!
					</button>
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
							<span class="text-gray-700 dark:text-gray-300">Local development</span>
						</div>

						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">Request logging</span>
						</div>
						<div class="flex items-center">
							<i class="fas fa-check text-green-500 mr-3"></i>
							<span class="text-gray-700 dark:text-gray-300">SSO integration</span>
						</div>
					</div>

					<button
						class="w-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 py-3 px-6 rounded-lg font-medium hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
						title="Download Community Edition"
						aria-label="Download Community Edition"
					>
						<i class="fas fa-download mr-2"></i>
						Download
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

			<!-- Feature Comparison Table -->
			<div class="mt-16">
				<div class="text-center mb-8">
					<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Compare Plans</h3>
					<p class="text-gray-600 dark:text-gray-400">See what's included in each plan</p>
				</div>

				<div
					class="bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden"
				>
					<div class="overflow-x-auto">
						<table class="w-full">
							<thead class="bg-gray-50 dark:bg-gray-700">
								<tr>
									<th class="px-6 py-4 text-left text-sm font-medium text-gray-900 dark:text-white"
										>Feature</th
									>
									<th
										class="px-6 py-4 text-center text-sm font-medium text-gray-900 dark:text-white"
										>Community</th
									>
									<th
										class="px-6 py-4 text-center text-sm font-medium text-gray-900 dark:text-white"
										>Cloud</th
									>
									<th
										class="px-6 py-4 text-center text-sm font-medium text-gray-900 dark:text-white"
										>Pro</th
									>
								</tr>
							</thead>
							<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
								<tr>
									<td class="px-6 py-4 text-sm text-gray-900 dark:text-white">Mock servers</td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
								</tr>
								<tr class="bg-gray-50 dark:bg-gray-750">
									<td class="px-6 py-4 text-sm text-gray-900 dark:text-white">Public hosting</td>
									<td class="px-6 py-4 text-center"><i class="fas fa-times text-red-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
								</tr>
								<tr>
									<td class="px-6 py-4 text-sm text-gray-900 dark:text-white">Team collaboration</td
									>
									<td class="px-6 py-4 text-center"><i class="fas fa-times text-red-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
								</tr>
								<tr class="bg-gray-50 dark:bg-gray-750">
									<td class="px-6 py-4 text-sm text-gray-900 dark:text-white">Advanced analytics</td
									>
									<td class="px-6 py-4 text-center"><i class="fas fa-times text-red-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-times text-red-500"></i></td>
									<td class="px-6 py-4 text-center"><i class="fas fa-check text-green-500"></i></td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</div>
		</div>
	</section>
</main>
