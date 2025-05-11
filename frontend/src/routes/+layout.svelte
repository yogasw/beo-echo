<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import ConfigurationList from '$lib/components/ConfigurationList.svelte';
	import Header from '$lib/components/Header.svelte';
	import { isOwnAuth, removeLocalStorage } from '$lib/utils/localStorage';
	import { getProjects } from '$lib/api/mockoonApi';
	import { onMount } from 'svelte';
	import { projects } from '$lib/stores/configurations';
	import Toast from '$lib/components/Toast.svelte';
	import { isAuthenticated } from '$lib/stores/authentication';
	import { browser } from '$app/environment';

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

	let searchTerm = '';
	let selectedConfig: Config | null = null;
	let activeTab = 'routes';
	// Check authentication from localStorage
	$: isLoginPage = $page.url.pathname === '/login';

	async function fetchConfigs() {
		try {
			await getProjects().then(d => {
				projects.set(d);
			});
		} catch (err) {
			console.error('Failed to fetch configs:', err);
		}
	}

	onMount(async () => {
		console.log("onMount: layout");
		if (isOwnAuth() && !$isAuthenticated && !isLoginPage) {
			await getProjects().then(async d => {
				isAuthenticated.set(true)
				await goto('/home');
			}).catch(async e => {
				console.error('Failed to fetch configs:', e);
				isAuthenticated.set(false)
				await goto('/login');
			})
		}

		async function initialize() {
			if ($isAuthenticated) {
				await fetchConfigs();
			}
			return () => {};
		}

		initialize();
	});

	function handleConfigSelect(event: CustomEvent<Config>) {
		selectedConfig = { ...event.detail };
	}

	function handleConfigStart(event: CustomEvent<Config>) {
		const config = event.detail;
		projects.update(configs => configs.map(c =>
			c.name === config.name ? { ...c, inUse: true } : c
		));
	}

	function handleConfigStop(event: CustomEvent<Config>) {
		const config = event.detail;
		projects.update(configs => configs.map(c =>
			c.name === config.name ? { ...c, inUse: false } : c
		));
	}

	function handleTabChange(event: CustomEvent<string>) {
		activeTab = event.detail;
		if (activeTab === 'routes') {
			goto('/home');
		} else if (activeTab === 'settings') {
			goto('/settings');
		}
	}

	function handleLogout() {
		isAuthenticated.set(false);
		removeLocalStorage('username');
		removeLocalStorage('password');
		goto("/login")
	}
</script>

{#if isLoginPage || !$isAuthenticated}
	<slot />
{:else}
	<div class="min-h-screen w-full bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-white font-sans">
		<div class="mmx-auto flex h-screen">
			<ConfigurationList
				{searchTerm}
				on:selectConfiguration={handleConfigSelect}
				on:startConfiguration={handleConfigStart}
				on:stopConfiguration={handleConfigStop}
			/>

			<div class="flex-1 flex flex-col overflow-hidden">
				<Header on:tabChange={handleTabChange} handleLogout={handleLogout} />
				<div class="flex-1 overflow-auto">
					<slot activeTab={activeTab} />
				</div>
			</div>

		</div>
	</div>

{/if}

<Toast />

<style global>

</style>
