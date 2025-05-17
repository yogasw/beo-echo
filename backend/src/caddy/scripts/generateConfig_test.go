package scripts

import (
	"context"
	"fmt"
	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/utils"
	"os"
	"testing"

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
	utils.SetupFolderConfigForTest()

	err := GenerateSingleConfigFromText(ctx, sampleConfigs)
	if err != nil {
		fmt.Println("Error generating config:", err)
	}
	require.NoError(t, err)

	// remove all generated files
	err = os.RemoveAll(lib.CONFIGS_DIR)
	require.NoError(t, err)
}
