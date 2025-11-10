<script lang="ts">
	// reference https://dev.to/lawrencecchen/monaco-editor-svelte-kit-572
	import { onMount, onDestroy, createEventDispatcher, afterUpdate } from 'svelte';
	import type * as monacoType from 'monaco-editor';
	import { theme as appTheme } from '$lib/stores/theme';

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
		});

		editor.onDidChangeModelContent(() => {
			if (!isUpdating) {
				currentValue = editor?.getValue() ?? '';
				dispatch('change', currentValue);
			}
		});
	});

	afterUpdate(() => {
		if (editor && value !== currentValue) {
			isUpdating = true;
			currentValue = value;
			editor.setValue(value);
			isUpdating = false;
		}
	});

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

<div bind:this={container} class="w-full h-full"></div>
