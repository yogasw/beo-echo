<script>
	function toggleUpdateRow(rowId) {
		const updateButton = document.getElementById(`update-button-${rowId}`);
		updateButton.classList.toggle("hidden");
	}

	function lockValue(rowId) {
		const valueInput = document.getElementById(`value-${rowId}`);
		valueInput.value = "xxx"; // Mask the value for security
		valueInput.disabled = true; // Lock the input
	}

	function showUpdateLoading(rowId) {
		const updateButton = document.getElementById(`update-button-${rowId}`);
		updateButton.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Updating...';
		updateButton.disabled = true;

		// Simulate a delay for the update process
		setTimeout(() => {
			updateButton.innerHTML = '<i class="fas fa-check"></i> Updated';
			updateButton.disabled = false;
		}, 3000); // 3 seconds delay
	}

</script>
<!-- Advanced Settings Tab -->
<div id="advancedSettingsTab" class="tab-content hidden">
	<div class="space-y-4">
		<h3 class="text-xl font-bold mb-4">Advanced Settings</h3>
		<div class="flex items-center mb-4">
			<div class="relative w-full">
				<i class="fas fa-search absolute top-3 left-3 text-gray-400"></i>
				<input
					type="text"
					class="w-full bg-gray-700 text-white py-2 pl-10 pr-4 rounded"
					placeholder="Search configuration"
				/>
			</div>
		</div>
		<hr class="border-gray-600 mb-4" />
		<div class="space-y-4">
			<!-- Row 1 -->
			<div class="flex items-center space-x-4">
				<input
					id="key-1"
					type="text"
					class="w-1/4 bg-gray-700 text-white py-2 px-4 rounded"
					placeholder="Key"
					on:click={()=>toggleUpdateRow(1)}
				/>
				<div class="relative w-3/4">
					<input
						id="value-1"
						type="text"
						class="w-full bg-gray-700 text-white py-2 px-4 rounded"
						placeholder="Single-line Value"
						on:click={()=>toggleUpdateRow(1)}
					/>
					<button
						class="absolute top-2 right-2 text-gray-400 hover:text-white"
						on:click={()=>lockValue(1)}
					>
						<i class="fas fa-lock"></i>
					</button>
				</div>
				<button
					id="update-button-1"
					class="text-gray-400 hover:text-blue-500"
					on:click={()=>showUpdateLoading(1)}
				>
					<i class="fas fa-check text-lg"></i> Update
				</button>
			</div>
			<!-- Row 2 -->
			<div class="flex items-center space-x-4">
				<input
					id="key-2"
					type="text"
					class="w-1/4 bg-gray-700 text-white py-2 px-4 rounded"
					placeholder="Key"
					on:click={()=>toggleUpdateRow(2)}
				/>
				<div class="relative w-3/4">
                    <textarea
											id="value-2"
											class="w-full bg-gray-700 text-white py-2 px-4 rounded resize-none"
											rows="2"
											placeholder="Multi-line Value"
											on:click={()=>toggleUpdateRow(2)}
										></textarea>
					<button
						class="absolute top-2 right-2 text-gray-400 hover:text-white"
						on:click={()=>lockValue(2)}
					>
						<i class="fas fa-lock"></i>
					</button>
				</div>
				<button
					id="update-button-2"
					class="text-gray-400 hover:text-blue-500"
					on:click={()=>showUpdateLoading(2)}
				>
					<i class="fas fa-check text-lg"></i> Update
				</button>
			</div>
			<!-- Row 3 -->
			<div class="flex items-center space-x-4">
				<input
					id="key-3"
					type="text"
					class="w-1/4 bg-gray-700 text-white py-2 px-4 rounded"
					placeholder="Key"
					on:click={()=>toggleUpdateRow(3)}
				/>
				<div class="relative w-3/4">
                    <textarea
											id="value-3"
											class="w-full bg-gray-700 text-white py-2 px-4 rounded resize-none"
											rows="4"
											placeholder="Long Text Area"
											on:click={()=>toggleUpdateRow(3)}
										></textarea>
					<button
						class="absolute top-2 right-2 text-gray-400 hover:text-white"
						on:click={()=>lockValue(3)}
					>
						<i class="fas fa-lock"></i>
					</button>
				</div>
				<button
					id="update-button-3"
					class="text-gray-400 hover:text-blue-500"
					on:click={()=>showUpdateLoading(3)}
				>
					<i class="fas fa-check text-lg"></i> Update
				</button>
			</div>
			<!-- Row 4 -->
			<div class="flex items-center space-x-4">
				<input
					id="key-4"
					type="text"
					class="w-1/4 bg-gray-700 text-white py-2 px-4 rounded"
					placeholder="Key"
					on:click={()=>toggleUpdateRow(4)}
				/>
				<div class="relative w-3/4">
					<input
						id="value-4"
						type="text"
						class="w-full bg-gray-700 text-white py-2 px-4 rounded"
						placeholder="Single-line Value"
						on:click={()=>toggleUpdateRow(4)}
					/>
					<button
						class="absolute top-2 right-2 text-gray-400 hover:text-white"
						on:click={()=>lockValue(4)}
					>
						<i class="fas fa-lock"></i>
					</button>
				</div>
				<button
					id="update-button-4"
					class="text-gray-400 hover:text-blue-500"
					on:click={()=>showUpdateLoading(4)}
				>
					<i class="fas fa-check text-lg"></i> Update
				</button>
			</div>
			<!-- Row 5 -->
			<div class="flex items-center space-x-4">
				<input
					id="key-5"
					type="text"
					class="w-1/4 bg-gray-700 text-white py-2 px-4 rounded"
					placeholder="Key"
					on:click={()=>toggleUpdateRow(5)}
				/>
				<div class="relative w-3/4">
                    <textarea
											id="value-5"
											class="w-full bg-gray-700 text-white py-2 px-4 rounded resize-none"
											rows="3"
											placeholder="Multi-line Value"
											on:click={()=>toggleUpdateRow(5)}
										></textarea>
					<button
						class="absolute top-2 right-2 text-gray-400 hover:text-white"
						on:click={()=>lockValue(5)}
					>
						<i class="fas fa-lock"></i>
					</button>
				</div>
				<button
					id="update-button-5"
					class="text-gray-400 hover:text-blue-500"
					on:click={()=>showUpdateLoading(5)}
				>
					<i class="fas fa-check text-lg"></i> Update
				</button>
			</div>
		</div>
	</div>
</div>
