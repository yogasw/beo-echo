<script lang="ts">
	/**
	 * Beautiful Beo Echo loading component inspired by the favicon design
	 * Supports light/dark themes with smooth animations
	 * 
	 * @prop {string} message - Loading message to display
	 * @prop {boolean} animated - Whether to show animations (default: true)
	 * @prop {string} size - Size variant: 'sm', 'md', 'lg' (default: 'md')
	 */
	export let message: string = 'Loading...';
	export let animated: boolean = true;
	export let size: 'sm' | 'md' | 'lg' = 'md';

	// Size configurations
	const sizeConfig = {
		sm: {
			container: 'w-16 h-16',
			icon: 'w-12 h-12',
			text: 'text-sm'
		},
		md: {
			container: 'w-24 h-24',
			icon: 'w-20 h-20',
			text: 'text-base'
		},
		lg: {
			container: 'w-32 h-32',
			icon: 'w-28 h-28',
			text: 'text-lg'
		}
	};

	$: config = sizeConfig[size];
</script>

<div class="flex flex-col items-center justify-center space-y-6">
	<!-- Animated Beo Echo Logo -->
	<div class="relative {config.container}">
		<!-- Main logo container with pulse effect -->
		<div class="absolute inset-0 {animated ? 'animate-pulse' : ''} opacity-20">
			<div class="w-full h-full rounded-full bg-gradient-to-br from-green-400 to-blue-500 dark:from-green-500 dark:to-blue-600"></div>
		</div>
		
		<!-- SVG Logo based on favicon - Complete bird design -->
		<div class="relative {config.icon} mx-auto {animated ? 'beo-spin' : ''}">
			<svg viewBox="0 0 612 823" class="w-full h-full drop-shadow-lg" xmlns="http://www.w3.org/2000/svg">
				<!-- Complete bird shape from favicon -->
				<!-- Main body/head -->
				<path d="M427.853 43.3665C442.4 39.3665 457.667 39.7398 472.6 40.4465C489.507 42.6998 507.16 46.7132 520.733 57.6465C508.8 61.4732 496.467 68.0998 490.667 79.7798C489.347 82.3798 487.8 84.8732 485.973 87.1665C474.733 77.9932 461.387 70.1532 446.52 69.7398C434.227 69.9532 421.693 74.6065 413.013 83.4732C400.613 95.1665 393.467 112.02 393.013 128.993C392.347 147.26 401.16 164.86 414.213 177.246C426.48 188.766 442.893 194.633 459.187 197.607C459.56 197.207 460.32 196.42 460.693 196.033C460.773 196.766 460.947 198.26 461.04 198.993C470.16 199.113 479.32 199.473 488.413 198.607C496.227 213.727 502.813 229.62 506.4 246.3C515.787 285.06 517.68 325.727 510.4 365.007C506.04 384.113 500.693 403.18 492.16 420.887C490.667 419.98 488.06 420.22 486.5 419.5C466.06 416.353 446.813 411.62 426.48 407.82C414.147 405.073 401.64 403.14 389.093 401.713C395.04 394.22 399.573 385.74 404.387 377.5C417.96 354.047 425.08 326.82 424.173 299.7C423.893 282.273 418.4 265.087 409.067 250.42C408.293 247.327 407.707 244.153 406.267 241.286C393.827 217.66 382.173 193.606 368.987 170.366C358.747 151.233 347.453 132.646 337.707 113.273C343.467 101.206 352.52 90.8598 361.613 81.1132C371.64 71.5798 382.36 62.5531 394.627 56.0065C405.187 50.5931 416.28 46.0998 427.853 43.3665Z" 
					  class="fill-green-500 dark:fill-green-400 transition-colors duration-300" />
				
				<!-- Beak area -->
				<path d="M490.667 79.7798C496.467 68.0998 508.8 61.4732 520.733 57.6465C521.293 57.9665 522.387 58.6065 522.933 58.9265C508.013 63.5931 492.52 73.2465 488.787 89.4998L488.253 89.0198C487.48 88.3798 486.72 87.7665 485.973 87.1665C487.8 84.8732 489.347 82.3798 490.667 79.7798Z" 
					  class="fill-orange-400 dark:fill-orange-300 transition-colors duration-300" />
				
				<!-- Upper beak -->
				<path d="M488.787 89.4999C492.52 73.2465 508.013 63.5932 522.933 58.9265C523.947 59.7665 525.013 60.6198 526.067 61.4465C525.333 61.6198 523.88 61.9665 523.16 62.1265C516.507 63.9932 509.627 66.1132 504.36 70.7798C501.213 73.5932 498.373 76.7132 495.547 79.8198C490.653 85.6598 490.747 93.7399 490.84 100.94C490.76 106.233 493.187 111.033 495.027 115.86L494.92 116.313L494.48 118.153C493.987 118.513 493.013 119.22 492.52 119.58C492.48 118.087 492.267 116.62 491.893 115.18C488.92 106.927 488.613 98.1799 488.787 89.4999Z" 
					  class="fill-red-500 dark:fill-red-400 transition-colors duration-300" />
				
				<!-- Eye area -->
				<path d="M438.88 81.6998C443.787 80.0198 449.213 80.5931 453.88 82.7131C459.96 85.9931 463.053 93.0732 462.973 99.7665C461.813 106.833 456.373 112.966 449.36 114.62C442.2 115.553 433.28 114.033 429.68 106.98C423.333 98.1932 429.147 84.9665 438.88 81.6998ZM441.453 88.2065C434.733 90.3798 432.107 98.9798 436.24 104.66C440.453 113.473 454.933 108.7 455.6 99.7665C456.907 92.1532 448.493 85.9532 441.453 88.2065Z" 
					  class="fill-yellow-500 dark:fill-yellow-400 transition-colors duration-300" />
				
				<!-- Eye pupil -->
				<path d="M441.453 88.2065C448.493 85.9532 456.907 92.1532 455.6 99.7665C454.933 108.7 440.453 113.473 436.24 104.66C432.107 98.9798 434.733 90.3798 441.453 88.2065ZM437.56 98.3798C436.747 103.367 443.813 107.993 447.72 104.367C445.053 102.5 441.547 99.6065 444.907 96.4599C447.36 98.2732 449.733 100.5 452.907 100.967C453 97.4465 452.533 93.5398 449.133 91.6332C444.307 87.9265 437.227 92.7532 437.56 98.3798Z" 
					  class="fill-gray-900 dark:fill-white transition-colors duration-300" />
				
				<!-- Body/chest area -->
				<path d="M311.613 203.66C310.893 211.566 311.68 219.526 310.96 227.433L310.827 228.713C315.253 225.526 320.08 228.686 324.4 230.353C332.507 227.18 341.413 227.633 349.92 228.38C372.387 230.46 393.853 246.193 400.267 268.18C405.293 285.406 400.613 304.086 390.96 318.806C376.533 341.646 354.853 358.393 334.733 375.9C296.453 407.06 259.347 441.313 213.547 461.22C199.107 467.3 183.707 473.593 167.707 471.38C161.893 471.326 156.667 467.566 150.813 468.5C159.667 476.313 170.707 481.166 181.987 484.313C209.227 491.7 233.987 489.38 261 481.5C274.5 483.5 368 426 369.553 422.407C378.42 414.633 382.013 409.613 390 401C395.947 393.507 399.573 385.74 404.387 377.5C417.96 354.046 425.08 326.82 424.173 299.7C423.893 282.273 418.4 265.086 409.067 250.42C401.267 244.02 395.027 235.793 386.427 230.38C381.76 227.673 376.707 225.753 371.827 223.5C364.76 220.193 356.853 219.513 349.173 219.26C336.907 219.166 324.467 220.7 313.133 225.62C309.827 218.833 312.08 210.926 311.613 203.66Z" 
					  class="fill-teal-600 dark:fill-teal-500 transition-colors duration-300" />
				
				<!-- Wing details -->
				<path d="M324.4 230.353C332.507 227.18 341.413 227.633 349.92 228.38C372.387 230.46 393.853 246.193 400.267 268.18C405.293 285.407 400.613 304.086 390.96 318.806C376.533 341.646 354.853 358.393 334.733 375.9C296.453 407.06 259.347 441.313 213.547 461.22C199.107 467.3 183.707 473.593 167.707 471.38C161.893 471.327 156.667 467.567 150.813 468.5C147.6 466.353 144.427 464.153 141.413 461.74C140.507 461.14 139.447 460.48 138.5 460C139.047 459.307 139.747 458.767 140.347 458.127C144.733 456.753 147.653 452.993 150.96 450.02C164.373 438.446 176.64 425.633 189.333 413.287C198.027 405.047 206.08 396.167 214.547 387.713C235.227 367.367 251.067 342.873 266.853 318.726C276.04 303.526 285.347 288.393 293.973 272.873C298.16 266.273 301.987 259.473 306.267 252.927C311.453 244.807 316.093 235.78 324.4 230.353Z" 
					  class="fill-emerald-500 dark:fill-emerald-400 transition-colors duration-300" />
				
				<!-- Tail feathers -->
				<path d="M167.907 540.967C176.173 533.847 183.547 525.74 192.16 519.02C192.64 525.526 194.013 531.926 194.88 538.393C196.707 550.326 197.08 562.42 198.827 574.38C201.427 590.753 202.707 607.3 205.053 623.713C206.093 629.94 205.027 636.486 207.48 642.473C195.627 666.86 181.827 690.406 165.813 712.313C152.227 731.633 136.907 749.847 119.173 765.513C112.733 770.9 106.067 776.113 98.6 780.02C94.5334 782.1 89.4667 783.153 85.2 781.087C79.4934 777.887 77.64 770.913 76.76 764.886C76.1333 754.206 79.0133 743.726 81.8267 733.54C85.48 719.953 90.8933 706.94 95.1733 693.566C103.507 669.526 112.413 645.686 120.653 621.62C124.333 610.433 128 599.233 132.627 588.406C135.707 580.993 139.04 573.633 141.08 565.847C150.16 557.713 158.693 548.98 167.907 540.967Z" 
					  class="fill-blue-700 dark:fill-blue-600 transition-colors duration-300" />
			</svg>
		</div>

		<!-- Rotating ring around logo -->
		{#if animated}
			<div class="absolute inset-0 border-2 border-transparent border-t-green-500 dark:border-t-green-400 border-r-blue-500 dark:border-r-blue-400 rounded-full animate-spin beo-ring"></div>
		{/if}
	</div>

	<!-- Loading message -->
	<div class="text-center space-y-2">
		<p class="theme-text-primary font-medium {config.text}">{message}</p>
		
		<!-- Animated dots -->
		{#if animated}
			<div class="flex items-center justify-center space-x-1">
				<div class="w-2 h-2 bg-green-500 dark:bg-green-400 rounded-full animate-bounce beo-dot-1"></div>
				<div class="w-2 h-2 bg-blue-500 dark:bg-blue-400 rounded-full animate-bounce beo-dot-2"></div>
				<div class="w-2 h-2 bg-purple-500 dark:bg-purple-400 rounded-full animate-bounce beo-dot-3"></div>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Custom animations for the Beo Echo loader */
	@keyframes beo-spin {
		0% {
			transform: rotate(0deg) scale(1);
		}
		50% {
			transform: rotate(180deg) scale(1.02);
		}
		100% {
			transform: rotate(360deg) scale(1);
		}
	}

	@keyframes beo-ring {
		0% {
			transform: rotate(0deg);
			border-width: 2px;
		}
		50% {
			border-width: 3px;
		}
		100% {
			transform: rotate(360deg);
			border-width: 2px;
		}
	}

	/* Apply animations with smoother timing */
	.beo-spin {
		animation: beo-spin 4s ease-in-out infinite;
	}

	.beo-ring {
		animation: beo-ring 3s linear infinite;
	}

	/* Staggered dot animations */
	.beo-dot-1 {
		animation-delay: 0s;
	}

	.beo-dot-2 {
		animation-delay: 0.2s;
	}

	.beo-dot-3 {
		animation-delay: 0.4s;
	}

	/* Smooth transitions for theme changes */
	svg path {
		transition: fill 0.3s ease;
	}

	/* Subtle glow effect in dark mode */
	:global(.dark) svg {
		filter: drop-shadow(0 0 8px rgba(34, 197, 94, 0.3));
	}

	/* Reduce motion for accessibility */
	@media (prefers-reduced-motion: reduce) {
		.beo-spin {
			animation: none;
		}
		
		.beo-ring {
			animation: none;
		}
		
		.animate-bounce {
			animation: none;
		}
		
		.animate-pulse {
			animation: none;
		}
	}
</style>
