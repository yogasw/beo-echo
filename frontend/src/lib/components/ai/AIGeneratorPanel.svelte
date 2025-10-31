<script lang="ts">
	import { generateContent } from '$lib/api/aiApi';
	import { toast } from '$lib/stores/toast';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let isOpen: boolean = false;
	export let initialContent: string = '';
	export let onSave: (content: string) => void = () => {};
	export let onClose: () => void = () => {};
	export let contentType: string = '';

	interface ChatMessage {
		role: 'user' | 'ai';
		content: string;
		timestamp: Date;
	}

	let chatMessages: ChatMessage[] = [];
	let userInput: string = '';
	let isGenerating: boolean = false;
	let currentContent: string = '';
	let chatContainer: HTMLDivElement;

	// Dragging functionality
	let isDragging = false;
	let dragStartX = 0;
	let dragStartY = 0;
	let panelX = 100;
	let panelY = 100;
	let currentX = 100;
	let currentY = 100;

	// Initialize content when panel opens
	$: if (isOpen && initialContent && !currentContent) {
		currentContent = initialContent;
		// Add welcome message
		if (chatMessages.length === 0) {
			chatMessages = [
				{
					role: 'ai',
					content: 'Hi! I can help you generate mock response data. Tell me what you need!',
					timestamp: new Date()
				}
			];
		}
	}

	function startDrag(e: MouseEvent) {
		if ((e.target as HTMLElement).classList.contains('drag-handle')) {
			isDragging = true;
			dragStartX = e.clientX - currentX;
			dragStartY = e.clientY - currentY;
		}
	}

	function drag(e: MouseEvent) {
		if (isDragging) {
			e.preventDefault();
			currentX = e.clientX - dragStartX;
			currentY = e.clientY - dragStartY;
			panelX = currentX;
			panelY = currentY;
		}
	}

	function stopDrag() {
		isDragging = false;
	}

	async function handleSendMessage() {
		if (!userInput.trim() || isGenerating) return;

		const message = userInput.trim();
		userInput = '';

		// Add user message
		chatMessages = [
			...chatMessages,
			{
				role: 'user',
				content: message,
				timestamp: new Date()
			}
		];

		// Scroll to bottom
		setTimeout(() => {
			if (chatContainer) {
				chatContainer.scrollTop = chatContainer.scrollHeight;
			}
		}, 0);

		try {
			isGenerating = true;

			const response = await generateContent({
				message: message,
				context: currentContent || undefined,
				content_type: contentType || undefined
			});

			let aiContent = response.content;
			let aiData = response.data;

			// Try to format as JSON if it's applicable
			if (response.can_apply) {
				try {
					const parsed = JSON.parse(aiData);
					aiData = JSON.stringify(parsed, null, 2);
				} catch (e) {
					// Not JSON, keep as is
				}
			}

			// Update current content
			currentContent = aiContent;

			// Auto-apply to editor only if can_apply is true
			if (response.can_apply) {
				onSave(aiData);
				toast.success('AI response applied to editor');

				// Add AI response
				chatMessages = [
					...chatMessages,
					{
						role: 'ai',
						content: aiContent,
						timestamp: new Date()
					},
					{
						role: 'ai',
						content: 'Generated and applied to editor!',
						timestamp: new Date()
					}
				];
			} else {
				// Just show the AI response without applying
				chatMessages = [
					...chatMessages,
					{
						role: 'ai',
						content: aiContent,
						timestamp: new Date()
					}
				];
			}
		} catch (error: any) {
			const errorMsg = error?.response?.data?.error || error?.message || 'Unknown error';

			chatMessages = [
				...chatMessages,
				{
					role: 'ai',
					content: `Error: ${errorMsg}`,
					timestamp: new Date()
				}
			];

			toast.error(`Failed: ${errorMsg}`);
		} finally {
			isGenerating = false;

			// Scroll to bottom
			setTimeout(() => {
				if (chatContainer) {
					chatContainer.scrollTop = chatContainer.scrollHeight;
				}
			}, 0);
		}
	}

	function handleKeyPress(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			handleSendMessage();
		}
	}

	function handleClose() {
		chatMessages = [];
		userInput = '';
		currentContent = '';
		onClose();
	}

	function clearChat() {
		chatMessages = [
			{
				role: 'ai',
				content: 'Chat cleared! How can I help you?',
				timestamp: new Date()
			}
		];
	}
</script>

<svelte:window on:mousemove={drag} on:mouseup={stopDrag} />

{#if isOpen}
	<div
		class="fixed {ThemeUtils.themeBgPrimary()} rounded-lg shadow-2xl border-2 {ThemeUtils.themeBorder()} z-50 flex flex-col"
		style="left: {panelX}px; top: {panelY}px; width: 450px; height: 600px;"
	>
		<!-- Header (Draggable) -->
		<div
			class="drag-handle flex items-center justify-between p-3 border-b {ThemeUtils.themeBorder()} cursor-move bg-gradient-to-r from-purple-600 to-pink-600 rounded-t-lg"
			on:mousedown={startDrag}
			role="button"
			tabindex="0"
		>
			<div class="flex items-center space-x-2">
				<i class="fas fa-magic text-white"></i>
				<h3 class="font-bold text-white text-sm">AI Assistant</h3>
			</div>
			<div class="flex items-center space-x-2">
				<button
					on:click={clearChat}
					class="text-white hover:text-gray-200 transition-colors text-xs"
					aria-label="Clear chat"
					title="Clear chat"
				>
					<i class="fas fa-trash"></i>
				</button>
				<button
					on:click={handleClose}
					class="text-white hover:text-gray-200 transition-colors"
					aria-label="Close panel"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>
		</div>

		<!-- Chat Messages -->
		<div
			bind:this={chatContainer}
			class="flex-grow overflow-y-auto p-4 space-y-3"
			style="max-height: calc(600px - 140px);"
		>
			{#each chatMessages as message}
				<div class="flex {message.role === 'user' ? 'justify-end' : 'justify-start'}">
					<div
						class="max-w-[85%] rounded-lg p-3 {message.role === 'user'
							? 'bg-blue-600 text-white'
							: ThemeUtils.themeBgSecondary() + ' ' + ThemeUtils.themeTextPrimary()}"
					>
						<div class="flex items-start space-x-2">
							<i
								class="fas {message.role === 'user' ? 'fa-user' : 'fa-robot'} mt-1 text-sm"
							></i>
							<div class="flex-grow">
								<p class="text-xs font-mono whitespace-pre-wrap break-words">
									{message.content}
								</p>
								<span class="text-xs opacity-60 mt-1 block">
									{message.timestamp.toLocaleTimeString()}
								</span>
							</div>
						</div>
					</div>
				</div>
			{/each}

			{#if isGenerating}
				<div class="flex justify-start">
					<div
						class="max-w-[85%] rounded-lg p-3 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()}"
					>
						<div class="flex items-center space-x-2">
							<i class="fas fa-robot text-sm"></i>
							<div class="flex space-x-1">
								<span class="animate-bounce">.</span>
								<span class="animate-bounce" style="animation-delay: 0.1s">.</span>
								<span class="animate-bounce" style="animation-delay: 0.2s">.</span>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>

		<!-- Input Area -->
		<div class="border-t {ThemeUtils.themeBorder()} p-3">
			<div class="flex space-x-2">
				<textarea
					bind:value={userInput}
					on:keypress={handleKeyPress}
					placeholder="Ask AI to generate data... (Enter to send, Shift+Enter for new line)"
					class="flex-grow px-3 py-2 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()} rounded-lg text-xs resize-none focus:ring-2 focus:ring-purple-500"
					rows="2"
					disabled={isGenerating}
				></textarea>
				<button
					on:click={handleSendMessage}
					disabled={!userInput.trim() || isGenerating}
					class="px-4 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg hover:from-purple-700 hover:to-pink-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center"
					aria-label="Send message"
				>
					<i class="fas fa-paper-plane"></i>
				</button>
			</div>
			<p class="text-xs {ThemeUtils.themeTextSecondary()} mt-2">
				<i class="fas fa-info-circle mr-1"></i>
				Drag header to move • Responses auto-apply to editor
			</p>
		</div>
	</div>
{/if}

<style>
	.drag-handle {
		user-select: none;
	}
</style>
