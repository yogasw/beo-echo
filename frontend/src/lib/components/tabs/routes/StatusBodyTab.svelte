<script lang="ts">
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import AIGeneratorButton from '$lib/components/ai/AIGeneratorButton.svelte';
	import { toast } from '$lib/stores/toast';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let responseBody: string;
	export let headers: string;
	export let statusCode: number;
	export let onStatusCodeChange: (val: number) => void;
	export let onSaveButtonClick: (body: string) => void;

	let editorRef: InstanceType<typeof MonacoEditor>;
	let isFullScreen = false;
	let contentType = '';

	//on mount, get content type from headers
	$: {
		
		try {
			const headersObj = JSON.parse(headers);
			contentType = headersObj['Content-Type'] || headersObj['content-type'] || '';
		} catch (e) {
			contentType = '';
		}
	}

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

	function handleAIGenerated(content: string) {
		responseBody = content;
		toast.success('AI generated content applied to editor');
	}
</script>

<div class="h-full flex flex-col space-y-2 w-full">
	<div>
		<label for="statusCode" class="text-sm {ThemeUtils.themeTextPrimary()}">Status Code:</label>
		<input
			id="statusCode"
			type="number"
			class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} p-2 rounded w-24 focus:outline-none focus:ring-0 focus:border-none"
			bind:value={statusCode}
			on:blur={(e) => {
				const input = e.target as HTMLInputElement;
				onStatusCodeChange(+input.value);
			}}
		/>
	</div>

	<div class="relative flex-grow w-full">
		<div class="absolute top-2 right-2 z-10 flex space-x-2">
			<AIGeneratorButton
				initialContent={responseBody}
				buttonText="AI Chat"
				size="sm"
				contentType={contentType}
				onGenerated={handleAIGenerated}
			/>
			<button
				on:click={formatContent}
				class="bg-green-600 text-white text-xs px-2 py-1 rounded hover:bg-green-700"
				aria-label="Format and prettify JSON content"
				title="Format and prettify JSON content"
			>
				Prettify
			</button>
			<button
				on:click={saveContent}
				class="bg-blue-600 text-white text-xs px-2 py-1 rounded hover:bg-blue-700"
				aria-label="Save response body content"
				title="Save response body content"
			>
				Save
			</button>
			<button
				on:click={toggleFullScreen}
				class="bg-gray-600 text-white text-xs px-2 py-1 rounded hover:bg-gray-700"
				aria-label="Toggle fullscreen editor"
				title="Toggle fullscreen editor"
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
			role="dialog"
			aria-label="Full Screen Editor"
			on:click={handleModalClick}
			on:keydown={(e) => {
				if (e.key === 'Escape') {
					toggleFullScreen();
				}
			}}
			tabindex="0"
		>
			<div class="relative w-11/12 h-5/6 {ThemeUtils.themeBgPrimary()} rounded-lg shadow-lg">
				<div class="absolute top-2 right-2 flex space-x-2 z-50">
					<button
						class="bg-blue-600 text-white text-xs px-3 py-1 rounded hover:bg-blue-700"
						on:click={saveContent}
						aria-label="Save fullscreen content"
						title="Save fullscreen content"
					>
						Save
					</button>
					<button
						class="bg-red-600 text-white text-xs px-3 py-1 rounded hover:bg-red-700"
						on:click={toggleFullScreen}
						aria-label="Close fullscreen editor"
						title="Close fullscreen editor"
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
