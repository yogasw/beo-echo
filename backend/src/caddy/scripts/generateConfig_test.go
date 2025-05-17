package scripts

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/utils"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// init db for test
	errDB := database.CheckAndHandle()
	assert.NoError(t, errDB)

	// Initialize main config
	errInit := InitCaddyConfig(ctx)
	if errInit != nil {
		fmt.Println("Error initializing Caddy config:", errInit)
	}
	assert.NoError(t, errInit)

	err := GenerateSingleConfigFromText(ctx, sampleConfigs)
	if err != nil {
		fmt.Println("Error generating config:", err)
	}
	require := assert.New(t)
	require.NoError(err)
}
