<script lang="ts">
	import type { Project } from '$lib/api/BeoApi';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { reconnectLogStream, refreshLogs, clearProjectLogs } from '$lib/services/logsService';

	export let selectedProject: Project;
	export let logsConnectionStatus: any;
	export let searchTerm: string = '';

</script>

{#if !logsConnectionStatus.isConnected && logsConnectionStatus.reconnectAttempts > 0}
	<div
		class="bg-red-100/30 dark:bg-red-900/30 border border-red-300 dark:border-red-700 p-2 rounded mb-4 flex items-center justify-between"
	>
		<div class="flex items-center">
			<i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 text-lg mr-2"></i>
			<span class="theme-text-primary">Live stream disconnected. Using manual refresh only.</span>
		</div>
		<button
			class={ThemeUtils.primaryButton('py-1 px-3 text-sm')}
			on:click={() => reconnectLogStream()}
			aria-label="Reconnect to live log stream"
			title="Reconnect to live log stream"
		>
			<i class="fas fa-sync mr-1"></i> Reconnect Stream
		</button>
	</div>
{/if}
<div class="mb-6">
	<div class="flex justify-between items-center mb-4">
		<div class="flex items-center">
			<div class="bg-blue-600/10 dark:bg-blue-600/10 p-2 rounded-lg mr-3">
				<i class="fas fa-list-alt text-blue-500 text-xl"></i>
			</div>
			<div>
				<h2 class="text-xl font-bold theme-text-primary">{selectedProject.name}</h2>
				<p class="text-sm theme-text-muted">Request logs</p>
			</div>
			<div class="ml-4 flex items-center bg-gray-100/50 dark:bg-gray-900/50 px-3 py-1 rounded-full">
				<!-- Stream status indicator -->
				<span class="relative flex h-3 w-3 mr-2">
					{#if logsConnectionStatus.isConnected}
						<span
							class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"
						></span>
						<span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
					{:else}
						<span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
					{/if}
				</span>
				<span
					class="text-xs font-medium {logsConnectionStatus.isConnected
						? 'text-green-400'
						: 'text-red-400'}"
				>
					{logsConnectionStatus.isConnected ? 'Live' : 'Offline'}
				</span>
			</div>
		</div>

		<div class="flex items-center space-x-3">
			<div class="flex items-center bg-gray-100/50 dark:bg-gray-900/50 px-3 py-1 rounded-full">
				<label class="inline-flex items-center cursor-pointer">
					<input
						type="checkbox"
						bind:checked={logsConnectionStatus.autoScroll}
						class="sr-only peer"
					/>
					<div
						class="relative w-9 h-5 bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-500 rounded-full peer peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all"
					></div>
				</label>
			</div>

			<div class="flex space-x-2">
				<button
					on:click={() => refreshLogs()}
					aria-label="Manually refresh request logs"
					title="Manually refresh request logs"
				>
					<i class="fas fa-sync mr-2"></i>
				</button>

				<button
					on:click={clearProjectLogs}
					aria-label="Clear all non-bookmarked logs"
					title="Clear all non-bookmarked logs (bookmarked logs will be preserved)"
				>
					<i class="fas fa-trash-alt mr-2"></i>
				</button>
			</div>
		</div>
	</div>

	<div class="relative mb-6">
		<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
			<i class="fas fa-search theme-text-muted"></i>
		</div>
		<input
			type="text"
			bind:value={searchTerm}
			placeholder="Search by keywords separated by spaces (e.g. 'GET users')..."
			class={ThemeUtils.inputField('p-3 ps-10 text-sm rounded-lg')}
		/>
	</div>
</div>
