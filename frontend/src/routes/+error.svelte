<script>
	export let error;
	export let status;

	let randomIcon = '❓';
	function getRandomEmoji() {
		// Common emoji range: Emoticons (U+1F600–U+1F64F)
		const min = 0x1f600;
		const max = 0x1f64f;
		const codePoint = Math.floor(Math.random() * (max - min + 1)) + min;
		return String.fromCodePoint(codePoint);
	}

	randomIcon = getRandomEmoji();
</script>

{#if status === 404}
	<div class="bg-gray-900 text-white h-screen flex items-center justify-center">
		<div class="text-center max-w-md">
			<div class="text-7xl mb-4">{randomIcon}</div>
			<h1 class="text-5xl font-bold mb-4">404</h1>
			<p class="text-xl text-gray-300 mb-6">
				Uh-oh! That page doesn’t exist. Here’s something to smile about:
			</p>
			<a href="/" class="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded inline-flex items-center">
				<i class="fas fa-arrow-left mr-2"></i>
				Return Home
			</a>
		</div>
	</div>
{:else}
	<h1>{status}</h1>
	<p>{error?.message}</p>
{/if}
