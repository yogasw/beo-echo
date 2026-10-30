package scripts

import (
	"context"
	"os"
	"testing"

	"beo-echo/backend/src/caddy/config"
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/lib"
	"beo-echo/backend/src/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCaddyConfig_MirrorSwitch(t *testing.T) {
	utils.SetupFolderConfigForTest()
	t.Cleanup(func() {
		lib.MIRROR_PROJECT_ALIAS = ""
		utils.CleanupTestFolders()
	})
	path := lib.CANDY_DIR + "/Caddyfile"

	// a normal Caddyfile is replaced once the instance becomes a mirror
	require.NoError(t, os.WriteFile(path, []byte(config.CADDY_DEFAULT_CONFIG), 0644))
	lib.MIRROR_PROJECT_ALIAS = "web-app"
	require.NoError(t, InitCaddyConfig(context.Background()))
	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, config.CADDY_MIRROR_CONFIG, string(got))

	// leaving mirror mode restores the default (UI) config; the subdomain
	// lookup that follows needs a DB (SQLite in tests)
	t.Setenv("DATABASE_URL", "")
	require.NoError(t, database.CheckAndHandle())
	lib.MIRROR_PROJECT_ALIAS = ""
	require.NoError(t, InitCaddyConfig(context.Background()))
	got, err = os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, config.CADDY_DEFAULT_CONFIG, string(got))
}
