package scripts

import (
	"beo-echo/backend/src/utils"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var sampleConfigs = []Config{
	{Domain: "*.test1.local", ProxyTarget: "localhost:8080"},
	{Domain: "*.test2.local", ProxyTarget: "localhost:8081"},
}

func TestGenerateSingleConfigFromText(t *testing.T) {
	ctx := context.Background()
	utils.SetupFolderConfigForTest()
	t.Cleanup(func() {
		utils.CleanupTestFolders()
	})

	err := GenerateSingleConfigFromText(ctx, sampleConfigs)
	if err != nil {
		fmt.Println("Error generating config:", err)
	}
	require.NoError(t, err)
}
