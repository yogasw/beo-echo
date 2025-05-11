<script lang="ts">
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import { updateResponse } from '$lib/stores/saveButton';
	import { toast } from '$lib/stores/toast';

	export let responseBody: string;
	export let statusCode: number;
	export let onStatusCodeChange: (val: number) => void;
	export let onSaveButtonClick: (body: string) => void;

	let editorRef: InstanceType<typeof MonacoEditor>;
	let isFullScreen = false;

	function formatContent() {
		editorRef?.format?.();
	}

	function saveContent() {
		const content: any = editorRef?.getValue();
		if (content) {
			onSaveButtonClick(content);
		}
	}

	function toggleFullScreen() {
		isFullScreen = !isFullScreen;
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape' && isFullScreen) {
			toggleFullScreen();
		}
	}

	function handleModalClick(event: MouseEvent) {
		if ((event.target as HTMLElement).classList.contains('fullscreen-modal')) {
			toggleFullScreen();
		}
	}
</script>

<div class="h-full flex flex-col space-y-2 w-full" on:keydown={handleKeyDown} tabindex="0">
	<div>
		<label class="text-sm text-white">Status Code:</label>
		<input
			type="number"
			class="bg-gray-700 text-white p-2 rounded w-24 focus:outline-none focus:ring-0 focus:border-none"
			bind:value={statusCode}
			on:blur={(e) => onStatusCodeChange(+e.target.value)}
		/>
	</div>

	<div class="relative flex-grow w-full">
		<div class="absolute top-2 right-2 z-10 flex space-x-2">
			<button
				on:click={formatContent}
				class="bg-green-600 text-white text-xs px-2 py-1 rounded hover:bg-green-700"
			>
				Prettify
			</button>
			<button
				on:click={saveContent}
				class="bg-blue-600 text-white text-xs px-2 py-1 rounded hover:bg-blue-700"
			>
				Save
			</button>
			<button
				on:click={toggleFullScreen}
				class="bg-gray-600 text-white text-xs px-2 py-1 rounded hover:bg-gray-700"
			>
				Full Screen
			</button>
		</div>
		<MonacoEditor
			bind:this={editorRef}
			value={responseBody}
			language="json"
		/>
	</div>

	{#if isFullScreen}
		<div
			class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-80 fullscreen-modal"
			on:click={handleModalClick}
		>
			<div class="relative w-11/12 h-5/6 bg-gray-800 rounded-lg shadow-lg">
				<div class="absolute top-2 right-2 flex space-x-2 z-50">
					<button
						class="bg-blue-600 text-white text-xs px-3 py-1 rounded hover:bg-blue-700"
						on:click={saveContent}
					>
						Save
					</button>
					<button
						class="bg-red-600 text-white text-xs px-3 py-1 rounded hover:bg-red-700"
						on:click={toggleFullScreen}
					>
						Close
					</button>
				</div>
				<MonacoEditor
					bind:this={editorRef}
					value={responseBody}
					language="json"
				/>
			</div>
		</div>
	{/if}
</div>
