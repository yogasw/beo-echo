package scripts

import (
	"context"
	"fmt"
	"mockoon-control-panel/backend_new/src/lib"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleConfigs = []Config{
	{Domain: "*.test1.local", ProxyTarget: "localhost:8080"},
	{Domain: "*.test2.local", ProxyTarget: "localhost:8081"},
}

func readFile(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}

func TestGenerateSingleConfigFromText(t *testing.T) {
	ctx := context.Background()

	// Ensure the directory exists when not creating a new config
	if _, err := os.Stat(lib.CANDY_DIR); os.IsNotExist(err) {
		err = os.MkdirAll(lib.CANDY_DIR, os.ModePerm)
		require.NoError(t, err)
	}
	// Clean up the directory before running the test
	if err := os.RemoveAll(lib.CANDY_DIR); err != nil {
		require.NoError(t, err)
	}
	defer func() {
		// Clean up after the test
		if err := os.RemoveAll(lib.CANDY_DIR); err != nil {
			require.NoError(t, err)
		}
	}()

	err := GenerateSingleConfigFromText(ctx, sampleConfigs)
	if err != nil {
		fmt.Println("Error generating config:", err)
	}

	require.NoError(t, err)

	dynContent := readFile(t, filepath.Join(lib.CANDY_DIR, "dynamic.conf"))
	assert.Contains(t, dynContent, "*.test1.local")
	assert.Contains(t, dynContent, "rewrite * /{http.request.host.labels.3}{http.request.uri}")
	dynamicImportPath := "./dynamic.conf"
	mainContent := readFile(t, filepath.Join(lib.CANDY_DIR, "Caddyfile"))
	assert.Contains(t, mainContent, "import "+dynamicImportPath)
}
