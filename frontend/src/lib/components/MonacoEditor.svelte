<script lang="ts">
	// reference https://dev.to/lawrencecchen/monaco-editor-svelte-kit-572
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';
	import type * as monacoType from 'monaco-editor';
	import { theme as appTheme } from '$lib/stores/theme';
	import { zoomLevel } from '$lib/stores/zoom';

	// Worker setup (modern)
	import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
	import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker';
	import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker';
	import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker';
	import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker';

	export let value: string = '{}';
	export let language: string = 'json';
	export let theme: string = 'vs-dark';

	let container: HTMLDivElement;
	let editor: monacoType.editor.IStandaloneCodeEditor | null = null;
	let monaco: typeof monacoType | null = null;

	const dispatch = createEventDispatcher();
	let currentValue = value;
	let isUpdating = false;

	// Monaco Environment for workers
	// @ts-ignore
	self.MonacoEnvironment = {
		getWorker: function (_: any, label: string) {
			if (label === 'json') return new jsonWorker();
			if (['css', 'scss', 'less'].includes(label)) return new cssWorker();
			if (['html', 'handlebars', 'razor'].includes(label)) return new htmlWorker();
			if (['typescript', 'javascript'].includes(label)) return new tsWorker();
			return new editorWorker();
		}
	};

	// Update editor theme when app theme changes
	$: if (monaco && editor && $appTheme) {
		const editorTheme = $appTheme === 'dark' ? 'vs-dark' : 'vs';
		monaco.editor.setTheme(editorTheme);
	}

	onMount(async () => {
		const m = await import('monaco-editor');
		monaco = m;
		
		// Set initial theme based on app theme
		const editorTheme = $appTheme === 'dark' ? 'vs-dark' : 'vs';
		
		editor = monaco.editor.create(container, {
			value,
			language,
			theme: editorTheme,
			automaticLayout: true,
			fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
			fontSize: 14,
			minimap: {
				enabled: false
			},
			scrollBeyondLastLine: false,
			fixedOverflowWidgets: true,
		});

		editor.onDidChangeModelContent(() => {
			if (!isUpdating) {
				currentValue = editor?.getValue() ?? '';
				dispatch('change', currentValue);
			}
		});
	});

	// Reactively sync value from parent
	$: updateContent(value);

	function updateContent(val: string) {
		if (editor && val !== currentValue) {
			isUpdating = true;
			currentValue = val;
			editor.setValue(val);
			isUpdating = false;
		}
	}

	onDestroy(() => {
		editor?.dispose();
		editor = null;
	});

	// Public Methods
	export function getValue(): string | undefined {
		return editor?.getValue();
	}

	export function setValue(val: string): void {
		if (editor) editor.setValue(val);
	}

	export function format(): void {
		editor?.getAction('editor.action.formatDocument')?.run();
	}

	export function setLanguage(lang: string): void {
		if (monaco && editor?.getModel()) {
			monaco.editor.setModelLanguage(editor.getModel()!, lang);
		}
	}
</script>

<!--
  Monaco Editor has a known issue where it does not support CSS `zoom` and its native coordinate system breaks.
  If the parent app applies a CSS zoom (e.g., 0.9), Monaco calculates mouse clicks with incorrect offsets, causing context menus to immediately close.
  To fix this, we apply an "inverted" zoom (1 / appZoom) specifically to Monaco's wrapper to force it back to a true 1.0 scale relative to the screen pixels.
  Then we manually adjust its width and height back down by multiplying by appZoom so it fits perfectly in its parent container.
-->
<div class="monaco-unzoom-wrapper" style="zoom: calc(1 / {$zoomLevel}); width: calc(100% * {$zoomLevel}); height: calc(100% * {$zoomLevel});">
	<div bind:this={container} class="w-full h-full relative"></div>
</div>
