package scripts

import (
	"context"
	"fmt"
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

func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "caddytest-*")
	require.NoError(t, err)
	return dir
}

// Copy original caddyfile to temp dir so test tidak overwrite file asli
func copyCaddyfileToTemp(t *testing.T, originalPath string) string {
	tmpDir := createTempDir(t)
	tempPath := filepath.Join(tmpDir, "Caddyfile")

	data, err := os.ReadFile(originalPath)
	require.NoError(t, err)

	err = os.WriteFile(tempPath, data, 0644)
	require.NoError(t, err)

	return tempPath
}

func readFile(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}

func TestGenerateSingleConfigFromText(t *testing.T) {
	ctx := context.Background()
	// Path asli Caddyfile dari struktur kamu
	originalCaddyfile := "../config/Caddyfile"

	// Salin dulu ke temp folder untuk test
	mainConfig := copyCaddyfileToTemp(t, originalCaddyfile)

	// Buat temp untuk Dynamic.conf
	tmpDir := filepath.Dir(mainConfig)
	dynamicConfig := filepath.Join(tmpDir, "Dynamic.conf")
	importPath := "./Dynamic.conf"

	err := GenerateSingleConfigFromText(ctx, dynamicConfig, mainConfig, importPath, sampleConfigs)
	if err != nil {
		fmt.Println("Error generating config:", err)
	}
	require.NoError(t, err)

	dynContent := readFile(t, dynamicConfig)
	assert.Contains(t, dynContent, "*.test1.local")
	assert.Contains(t, dynContent, "rewrite * /{http.request.host.labels.3}{http.request.uri}")

	mainContent := readFile(t, mainConfig)
	assert.Contains(t, mainContent, "import "+importPath)
}
