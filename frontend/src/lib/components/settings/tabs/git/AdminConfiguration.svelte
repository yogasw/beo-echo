<script>
	function togglePasswordVisibility() {
		const oldPasswordInput = document.getElementById("oldPassword");
		const newPasswordInput = document.getElementById("newPassword");
		const isChecked = document.getElementById("showPasswordCheckbox").checked;

		oldPasswordInput.type = isChecked ? "text" : "password";
		newPasswordInput.type = isChecked ? "text" : "password";
	}
	function showLoading(buttonId) {
		const button = document.getElementById(buttonId);
		button.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Syncing...';
		button.disabled = true;

		// Simulate a delay for the sync process
		setTimeout(() => {
			button.innerHTML = 'Save and Test Sync';
			button.disabled = false;
		}, 5000); // 5 seconds delay
	}

</script>

<!-- Admin Settings Tab -->
<div id="adminSettingsTab" class="tab-content">
	<div class="space-y-4">
		<div>
			<label class="block text-sm font-bold mb-2">Old Password</label>
			<input
				id="oldPassword"
				type="password"
				class="w-full bg-gray-700 text-white py-2 px-4 rounded"
				placeholder="Old password"
			/>
		</div>
		<div>
			<label class="block text-sm font-bold mb-2">New Password</label>
			<input
				id="newPassword"
				type="password"
				class="w-full bg-gray-700 text-white py-2 px-4 rounded"
				placeholder="New password"
			/>
		</div>
		<div class="flex items-center">
			<input
				id="showPasswordCheckbox"
				type="checkbox"
				class="mr-2"
				on:click={()=>togglePasswordVisibility()}
			/>
			<label for="showPasswordCheckbox" class="text-sm">
				Show Password
			</label>
		</div>
	</div>
	<button class="bg-green-500 text-white py-2 px-4 rounded mt-4 w-full">
		Update
	</button>
</div>
