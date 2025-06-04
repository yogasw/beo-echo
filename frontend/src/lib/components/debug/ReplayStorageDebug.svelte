<!-- ReplayStorageDebug component for testing project-based localStorage integration -->
<script lang="ts">
	import { 
		getReplayEditorState, 
		setReplayEditorState, 
		clearReplayEditorState,
		createDefaultReplayEditorState,
		cleanupProjectStorage,
		cleanupWorkspaceStorage,
		getStorageInfo,
		getStoredProjectIds
	} from '$lib/utils/replayEditorStorage';
	
	let testWorkspaceId = 'test-workspace-1';
	let testProjectId = 'test-project-1';
	let testResult = '';
	let storageInfo = '';

	// Test basic project-based storage
	function testProjectStorage() {
		try {
			// Create test state
			const state = createDefaultReplayEditorState(testWorkspaceId, testProjectId);
			state.tabs.push({
				id: 'test-tab-1',
				name: 'Test Tab',
				method: 'GET',
				url: 'https://api.example.com/test',
				isUnsaved: true
			});

			// Save state
			setReplayEditorState(state);
			
			// Load state back
			const loadedState = getReplayEditorState(testWorkspaceId, testProjectId);
			
			if (loadedState && loadedState.tabs.length === 2) {
				testResult = '✅ Project storage test PASSED';
			} else {
				testResult = '❌ Project storage test FAILED';
			}
		} catch (error) {
			testResult = `❌ Project storage test ERROR: ${error}`;
		}
		
		updateStorageInfo();
	}

	// Test project isolation
	function testProjectIsolation() {
		try {
			const workspace1 = 'workspace-1';
			const workspace2 = 'workspace-2';
			const project1 = 'project-1';
			const project2 = 'project-2';

			// Create different states for different projects
			const state1 = createDefaultReplayEditorState(workspace1, project1);
			state1.tabs[0].name = 'Project 1 Tab';
			
			const state2 = createDefaultReplayEditorState(workspace2, project2);
			state2.tabs[0].name = 'Project 2 Tab';

			// Save both states
			setReplayEditorState(state1);
			setReplayEditorState(state2);

			// Load and verify isolation
			const loaded1 = getReplayEditorState(workspace1, project1);
			const loaded2 = getReplayEditorState(workspace2, project2);

			if (loaded1?.tabs[0].name === 'Project 1 Tab' && 
				loaded2?.tabs[0].name === 'Project 2 Tab') {
				testResult = '✅ Project isolation test PASSED';
			} else {
				testResult = '❌ Project isolation test FAILED';
			}

			// Cleanup test data
			clearReplayEditorState(workspace1, project1);
			clearReplayEditorState(workspace2, project2);
		} catch (error) {
			testResult = `❌ Project isolation test ERROR: ${error}`;
		}
		
		updateStorageInfo();
	}

	// Test cleanup functions
	function testCleanup() {
		try {
			const workspace = 'cleanup-workspace';
			const project1 = 'cleanup-project-1';
			const project2 = 'cleanup-project-2';

			// Create test data
			setReplayEditorState(createDefaultReplayEditorState(workspace, project1));
			setReplayEditorState(createDefaultReplayEditorState(workspace, project2));

			// Verify data exists
			const storedProjects = getStoredProjectIds(workspace);
			if (storedProjects.length !== 2) {
				testResult = '❌ Cleanup test FAILED - setup issue';
				return;
			}

			// Test project cleanup
			cleanupProjectStorage(workspace, project1);
			const afterProjectCleanup = getStoredProjectIds(workspace);
			
			if (afterProjectCleanup.length !== 1) {
				testResult = '❌ Cleanup test FAILED - project cleanup failed';
				return;
			}

			// Test workspace cleanup
			cleanupWorkspaceStorage(workspace);
			const afterWorkspaceCleanup = getStoredProjectIds(workspace);
			
			if (afterWorkspaceCleanup.length === 0) {
				testResult = '✅ Cleanup test PASSED';
			} else {
				testResult = '❌ Cleanup test FAILED - workspace cleanup failed';
			}
		} catch (error) {
			testResult = `❌ Cleanup test ERROR: ${error}`;
		}
		
		updateStorageInfo();
	}

	function updateStorageInfo() {
		const info = getStorageInfo();
		const projects = getStoredProjectIds(testWorkspaceId);
		storageInfo = `Storage Info:
- Total keys: ${info.totalKeys}
- Replay editor keys: ${info.replayEditorKeys}
- Storage size: ${info.totalSizeKB.toFixed(2)} KB
- Test workspace projects: ${projects.length}`;
	}

	function clearAllTestData() {
		try {
			// Clear any test data
			cleanupWorkspaceStorage('test-workspace-1');
			cleanupWorkspaceStorage('workspace-1');
			cleanupWorkspaceStorage('workspace-2');
			cleanupWorkspaceStorage('cleanup-workspace');
			
			testResult = '🧹 All test data cleared';
			updateStorageInfo();
		} catch (error) {
			testResult = `❌ Clear test data ERROR: ${error}`;
		}
	}

	// Initialize storage info
	updateStorageInfo();
</script>

<div class="p-4 theme-bg-primary border theme-border rounded-lg">
	<h3 class="text-lg font-semibold theme-text-primary mb-4">
		<i class="fas fa-debug mr-2"></i>
		Replay Storage Debug
	</h3>

	<div class="space-y-4">
		<!-- Test Controls -->
		<div class="flex flex-wrap gap-2">
			<button 
				class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded text-sm"
				on:click={testProjectStorage}
				title="Test basic project-based storage functionality"
				aria-label="Test project storage"
			>
				<i class="fas fa-play mr-1"></i>
				Test Project Storage
			</button>
			
			<button 
				class="px-3 py-2 bg-green-600 hover:bg-green-700 text-white rounded text-sm"
				on:click={testProjectIsolation}
				title="Test that projects have isolated storage"
				aria-label="Test project isolation"
			>
				<i class="fas fa-shield-alt mr-1"></i>
				Test Isolation
			</button>
			
			<button 
				class="px-3 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded text-sm"
				on:click={testCleanup}
				title="Test cleanup functions for projects and workspaces"
				aria-label="Test cleanup functions"
			>
				<i class="fas fa-broom mr-1"></i>
				Test Cleanup
			</button>
			
			<button 
				class="px-3 py-2 bg-red-600 hover:bg-red-700 text-white rounded text-sm"
				on:click={clearAllTestData}
				title="Clear all test data from localStorage"
				aria-label="Clear test data"
			>
				<i class="fas fa-trash mr-1"></i>
				Clear Test Data
			</button>
		</div>

		<!-- Test Configuration -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div>
				<label for="workspace-id" class="block text-sm font-medium theme-text-secondary mb-1">
					Test Workspace ID
				</label>
				<input 
					id="workspace-id"
					type="text" 
					bind:value={testWorkspaceId}
					class="w-full p-2 theme-bg-secondary border theme-border rounded text-sm theme-text-primary"
					placeholder="Enter workspace ID"
				/>
			</div>
			
			<div>
				<label for="project-id" class="block text-sm font-medium theme-text-secondary mb-1">
					Test Project ID
				</label>
				<input 
					id="project-id"
					type="text" 
					bind:value={testProjectId}
					class="w-full p-2 theme-bg-secondary border theme-border rounded text-sm theme-text-primary"
					placeholder="Enter project ID"
				/>
			</div>
		</div>

		<!-- Test Results -->
		{#if testResult}
			<div class="p-3 theme-bg-secondary border theme-border rounded">
				<h4 class="font-medium theme-text-primary mb-2">Test Result:</h4>
				<pre class="text-sm theme-text-secondary whitespace-pre-wrap">{testResult}</pre>
			</div>
		{/if}

		<!-- Storage Info -->
		{#if storageInfo}
			<div class="p-3 theme-bg-secondary border theme-border rounded">
				<h4 class="font-medium theme-text-primary mb-2">Storage Information:</h4>
				<pre class="text-sm theme-text-secondary whitespace-pre-wrap">{storageInfo}</pre>
			</div>
		{/if}
	</div>
</div>
