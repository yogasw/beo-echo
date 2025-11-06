<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { ReplaceTextConfig } from '$lib/types/Action';

	export let config: ReplaceTextConfig | null = null;

	const dispatch = createEventDispatcher();

	// Initialize config with default values
	let target: 'request_body' | 'response_body' | 'request_header' | 'response_header' =
		config?.target || 'response_body';
	let pattern = config?.pattern || '';
	let replacement = config?.replacement || '';
	let useRegex = config?.use_regex ?? false;
	let headerKey = config?.header_key || '';

	// Update parent when any value changes
	$: {
		const newConfig: ReplaceTextConfig = {
			target,
			pattern,
			replacement,
			use_regex: useRegex,
			...(target.includes('header') && headerKey ? { header_key: headerKey } : {})
		};
		dispatch('change', newConfig);
	}

	// Show/hide header key field based on target
	$: showHeaderKey = target === 'request_header' || target === 'response_header';
</script>

<div class="space-y-4">
	<!-- Target Selection -->
	<div>
		<label for="target" class="block text-sm font-medium theme-text-primary mb-2">
			Target <span class="text-red-500">*</span>
		</label>
		<select
			id="target"
			bind:value={target}
			class="block w-full p-3 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
			required
			aria-required="true"
			aria-describedby="target-help"
		>
			<option value="response_body">Response Body</option>
			<option value="request_body">Request Body</option>
			<option value="response_header">Response Header</option>
			<option value="request_header">Request Header</option>
		</select>
		<p id="target-help" class="mt-1 text-xs theme-text-secondary">
			Choose where to perform the text replacement2
		</p>
	</div>

	<!-- Header Key (conditional) -->
	{#if showHeaderKey}
		<div>
			<label for="header-key" class="block text-sm font-medium theme-text-primary mb-2">
				Header Key <span class="text-red-500">*</span>
			</label>
			<input
				id="header-key"
				type="text"
				bind:value={headerKey}
				class="block w-full p-3 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400"
				placeholder="e.g., Content-Type, Authorization"
				required={showHeaderKey}
				aria-required={showHeaderKey}
				aria-describedby="header-key-help"
			/>
			<p id="header-key-help" class="mt-1 text-xs theme-text-secondary">
				The name of the header to modify
			</p>
		</div>
	{/if}

	<!-- Pattern -->
	<div>
		<label for="pattern" class="block text-sm font-medium theme-text-primary mb-2">
			Find Pattern <span class="text-red-500">*</span>
		</label>
		<input
			id="pattern"
			type="text"
			bind:value={pattern}
			class="block w-full p-3 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400 font-mono"
			placeholder={useRegex ? 'e.g., \\d{3}-\\d{4}' : 'e.g., localhost:3000'}
			required
			aria-required="true"
			aria-describedby="pattern-help"
		/>
		<p id="pattern-help" class="mt-1 text-xs theme-text-secondary">
			{useRegex
				? 'Enter a regular expression pattern to match'
				: 'Enter the exact text to find'}
		</p>
	</div>

	<!-- Replacement -->
	<div>
		<label for="replacement" class="block text-sm font-medium theme-text-primary mb-2">
			Replacement Text <span class="text-red-500">*</span>
		</label>
		<input
			id="replacement"
			type="text"
			bind:value={replacement}
			class="block w-full p-3 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400 font-mono"
			placeholder={useRegex ? 'e.g., $1-XXXX (use $1, $2 for capture groups)' : 'e.g., api.production.com'}
			required
			aria-required="true"
			aria-describedby="replacement-help"
		/>
		<p id="replacement-help" class="mt-1 text-xs theme-text-secondary">
			{useRegex
				? 'Text to replace with (supports $1, $2 for regex capture groups)'
				: 'Text to replace the pattern with'}
		</p>
	</div>

	<!-- Use Regex Toggle -->
	<div class="flex items-center space-x-3">
		<label class="relative inline-flex items-center cursor-pointer">
			<input
				id="use-regex"
				type="checkbox"
				bind:checked={useRegex}
				class="sr-only peer"
				aria-label="Use regular expression matching"
			/>
			<div
				class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
			></div>
		</label>
		<div>
			<label for="use-regex" class="text-sm font-medium theme-text-primary cursor-pointer">
				Use Regular Expression
			</label>
			<p class="text-xs theme-text-secondary">
				Enable regex for advanced pattern matching
			</p>
		</div>
	</div>

	<!-- Example Preview -->
	<div class="mt-4 p-3 bg-gray-900/50 rounded border border-blue-900/50">
		<div class="text-xs theme-text-secondary mb-2 flex items-center">
			<i class="fas fa-info-circle mr-2 text-blue-500"></i>
			Example Preview
		</div>
		<div class="space-y-2 text-xs">
			<div>
				<span class="theme-text-secondary">Input:</span>
				<code class="ml-2 px-2 py-1 bg-gray-800 rounded theme-text-primary font-mono">
					{pattern || '(text to find)'}
				</code>
			</div>
			<div class="flex items-center theme-text-secondary">
				<i class="fas fa-arrow-down mr-2"></i>
				Replaced with
			</div>
			<div>
				<span class="theme-text-secondary">Output:</span>
				<code class="ml-2 px-2 py-1 bg-gray-800 rounded text-green-400 font-mono">
					{replacement || '(replacement text)'}
				</code>
			</div>
		</div>
	</div>
</div>
