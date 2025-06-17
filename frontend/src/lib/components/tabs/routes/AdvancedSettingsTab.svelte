<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let delayMs: number = 0;
	export let onDelayUpdate: (delayMs: number) => void = () => {};

	// Update delay value and trigger callback
	function updateDelay(newDelayMs: number): void {
		// Validate the delay value
		const validatedDelay = validateDelay(newDelayMs);
		
		// Trigger the callback with the new delay value
		onDelayUpdate(validatedDelay);
	}

	// Handle delay reset
	function handleDelayReset() {
		updateDelay(0);
	}

	// Validate delay value
	function validateDelay(value: number): number {
		if (value < 0) return 0;
		if (value > 120000) return 120000;
		return value;
	}

	$: delaySeconds = delayMs / 1000;
</script>

<div class="space-y-4 h-full">
	<!-- Delay Configuration -->
	<div class="{ThemeUtils.card()}">
		<div class="p-4 border-b {ThemeUtils.themeBorder()} {ThemeUtils.themeBgSecondary()}">
			<h3 class="text-sm font-medium {ThemeUtils.themeTextPrimary()} flex items-center">
				<i class="fas fa-hourglass-half text-orange-500 mr-2"></i>
				Response Delay
			</h3>
		</div>
		<div class="p-4">
			<div class="flex items-start justify-between mb-4">
				<div class="flex-1">
					<label for="response-delay-ms" class="text-sm font-medium {ThemeUtils.themeTextSecondary()} block mb-1">
						Delay (milliseconds)
					</label>
					<p class="text-xs {ThemeUtils.themeTextMuted()}">
						Add artificial delay to this response. Set to 0 to disable delay.
					</p>
				</div>
				<button
					type="button"
					class="ml-4 px-3 py-1 text-xs bg-red-600 hover:bg-red-700 text-white rounded-md flex items-center"
					on:click={handleDelayReset}
					title="Reset delay to 0"
					aria-label="Reset delay configuration"
				>
					<i class="fas fa-undo mr-1"></i>
					Reset
				</button>
			</div>
			
			<div class="relative mb-2">
				<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
					<i class="fas fa-clock {ThemeUtils.themeTextMuted()}"></i>
				</div>
				<input
					type="number"
					id="response-delay-ms"
					min="0"
					max="120000"
					step="1000"
					class="{ThemeUtils.inputField()}"
					value={delayMs}
					on:input={(e) => {
						const value = parseInt(e.currentTarget.value);
						if (!isNaN(value)) {
							const validatedValue = validateDelay(value);
							updateDelay(validatedValue);
							if (validatedValue !== value) {
								e.currentTarget.value = validatedValue.toString();
							}
						}
					}}
					placeholder="Enter delay in milliseconds (0 = disabled)"
					title="Response delay in milliseconds (0-120000ms / 2 minutes max)"
					aria-label="Response delay in milliseconds"
				/>
			</div>
			<p class="text-xs {ThemeUtils.themeTextMuted()}">
				Current: {delaySeconds}s (Max: 120 seconds)
			</p>
		</div>
	</div>
</div>
