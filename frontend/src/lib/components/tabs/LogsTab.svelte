<script lang="ts">
	import {onDestroy, onMount } from 'svelte';
	import type { Project, RequestLog } from '$lib/api/BeoApi';
	import { addBookmark, deleteBookmark } from '$lib/api/BeoApi';
	import { fade } from 'svelte/transition';
	import ModalCreateMock from './logs/ModalCreateMock.svelte';
	import LogsHeader from './logs/LogsHeader.svelte';
	import LogsList from './logs/LogsList.svelte';
	import SearchNoResults from './logs/SearchNoResults.svelte';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { logs, logsConnectionStatus } from '$lib/stores/logs';
	import {reconnectLogStream, refreshLogs } from '$lib/services/logsService';

	export let selectedProject: Project;

	let searchTerm = '';
	const pageSize = 100;
	// Map to track expanded logs
	let expandedLogs: Record<string, boolean> = {};
	// Map to track active tab (request/response)
	let activeTabs: Record<string, 'request' | 'response'> = {};
	// Notification for copy operations
	let copyNotification = { show: false, message: '' };
	// State for create mock modal
	let isCreateMockModalOpen = false;
	let selectedLogForMock: RequestLog | null = null;

	// Function to toggle log expansion
	function toggleLogExpansion(logId: string) {
		expandedLogs[logId] = !expandedLogs[logId];
		if (expandedLogs[logId] && !activeTabs[logId]) {
			activeTabs[logId] = lastActiveTab; // Use the last active tab instead of hardcoded 'request'
		}
		expandedLogs = expandedLogs; // Force Svelte reactivity update
		activeTabs = activeTabs; // Force Svelte reactivity update
	}

	// Function to switch between request and response tabs
	let lastActiveTab: 'request' | 'response' = 'request'; // Store the last active tab
	
	function switchTab(logId: string, tab: 'request' | 'response') {
		activeTabs[logId] = tab;
		lastActiveTab = tab; // Remember this tab as the last used one
		activeTabs = activeTabs; // Force Svelte reactivity update
	}
	
	// Listen for the clearSearch event
	onMount(() => {
		const handleClearSearch = () => {
			searchTerm = '';
		};
		
		document.addEventListener('clearSearch', handleClearSearch);
		
		return () => {
			document.removeEventListener('clearSearch', handleClearSearch);
		};
	});

	// Function to copy content to clipboard
	async function copyToClipboard(text: string, label: string) {
		try {
			await navigator.clipboard.writeText(text);
			copyNotification = { show: true, message: `${label} copied to clipboard!` };
			setTimeout(() => {
				copyNotification = { show: false, message: '' };
			}, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
			copyNotification = { show: true, message: 'Failed to copy to clipboard' };
			setTimeout(() => {
				copyNotification = { show: false, message: '' };
			}, 2000);
		}
	}

	// Function to pretty format JSON
	function formatJson(jsonStr: string): string {
		try {
			return JSON.stringify(JSON.parse(jsonStr), null, 2);
		} catch (e) {
			return jsonStr;
		}
	}

	// Function to minify JSON
	function minifyJson(jsonStr: string): string {
		try {
			return JSON.stringify(JSON.parse(jsonStr));
		} catch (e) {
			return jsonStr;
		}
	}

	// Function to check if all search terms are present in a log
	function matchesAllSearchTerms(log: RequestLog, searchTerms: string[]): boolean {
		if (searchTerms.length === 0) return true;

		// Combine all searchable fields into one string for easier searching
		const searchableText = [
			log.path.toLowerCase(),
			log.method.toLowerCase(),
			log.request_body.toLowerCase(),
			log.response_body.toLowerCase()
		].join(' ');

		// Check if all search terms are present in the searchable text
		return searchTerms.every((term) => searchableText.includes(term));
	}

	$: searchTerms = searchTerm
		.toLowerCase()
		.split(' ')
		.filter((term) => term.trim() !== '');
	$: filteredLogs = searchTerm
		? $logs.filter((log) => matchesAllSearchTerms(log, searchTerms))
		: $logs;
	
	// Update auto-scroll setting in store when changed
	$: {
		// When autoScroll changes from the UI, update the store
		if ($logsConnectionStatus.autoScroll !== $logsConnectionStatus.autoScroll) {
			$logsConnectionStatus.autoScroll = $logsConnectionStatus.autoScroll;
		}
	}

	// Convert JSON string to object for display
	function parseJson(jsonString: string): any {
		try {
			return JSON.parse(jsonString);
		} catch (e) {
			return jsonString;
		}
	}

	// Format timestamp for display
	function formatDate(dateString: string | Date): string {
		try {
			const date = typeof dateString === 'string' ? new Date(dateString) : dateString;
			return date.toLocaleString();
		} catch (e) {
			return String(dateString);
		}
	}

	// Function to create a mock from a log entry
	function createMockFromLog(log: RequestLog) {
		// Open the modal for creating a mock from the log
		selectedLogForMock = log;
		isCreateMockModalOpen = true;
		console.log('Opening mock creation modal for log:', log.id);
	}
	
	// Function to bookmark a log
	async function bookmarkLog(log: RequestLog) {
		try {
			if (log.bookmark) {
				// If already bookmarked, remove the bookmark
				await deleteBookmark(selectedProject.id, log.id);
				log.bookmark = false;
				
				copyNotification = { show: true, message: 'Bookmark removed successfully!' };
			} else {
				// Otherwise add a bookmark
				await addBookmark(selectedProject.id, log);
				// The log object is updated by the API function to set bookmark=true
				
				copyNotification = { show: true, message: 'Log bookmarked successfully!' };
			}
			
			// Force Svelte reactivity to update the UI
			logs.update(logs => [...logs]);
			
			setTimeout(() => {
				copyNotification = { show: false, message: '' };
			}, 2000);
		} catch (err) {
			console.error('Failed to update bookmark:', err);
			copyNotification = { show: true, message: 'Failed to update bookmark' };
			setTimeout(() => {
				copyNotification = { show: false, message: '' };
			}, 2000);
		}
	}

	// Handle success after mock creation
	function handleMockCreationSuccess() {
		copyNotification = { show: true, message: 'Mock endpoint created successfully!' };
		setTimeout(() => {
			copyNotification = { show: false, message: '' };
		}, 2000);
	}

	// Clean up on component destroy
	onDestroy(() => {
		// No cleanup needed as the closeLogStream is handled at the service level
	});
</script>

<div class="w-full theme-bg-primary p-4 relative">
	<!-- Copy notification toast -->
	{#if copyNotification.show}
		<div
			transition:fade={{ duration: 200 }}
			class="fixed top-6 right-6 theme-bg-secondary theme-text-primary px-4 py-2 rounded shadow-lg z-50 flex items-center"
		>
			<i class="fas fa-check-circle text-green-400 mr-2"></i>
			<span>{copyNotification.message}</span>
		</div>
	{/if}

	<LogsHeader 
		{selectedProject} 
		logsConnectionStatus={$logsConnectionStatus} 
		bind:searchTerm={searchTerm} 
	/>

	{#if $logsConnectionStatus.isLoading}
		<div class="flex justify-center py-8">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if $logsConnectionStatus.error}
		<div class="bg-red-100 dark:bg-red-800 p-4 rounded mb-4 text-center">
			<p class="text-red-700 dark:text-white">{$logsConnectionStatus.error}</p>
			<button on:click={() => refreshLogs()} class={ThemeUtils.primaryButton('mt-2 py-1 px-4 text-sm')}>
				Retry
			</button>
		</div>
	{:else if filteredLogs.length === 0 && searchTerm}
		<SearchNoResults {searchTerm} />
	{:else}
		<LogsList 
			{selectedProject}
			filteredLogs={filteredLogs}
			{expandedLogs}
			{activeTabs}
			logsConnectionStatus={$logsConnectionStatus}
			{toggleLogExpansion}
			{switchTab}
			{copyToClipboard}
			{parseJson}
			{formatDate}
			{bookmarkLog}
			{createMockFromLog}
		/>
	{/if}

	<!-- Create Mock Modal -->
	<ModalCreateMock
		isOpen={isCreateMockModalOpen}
		log={selectedLogForMock}
		projectId={selectedProject.id}
		onClose={() => (isCreateMockModalOpen = false)}
		onSuccess={handleMockCreationSuccess}
	/>
</div>
