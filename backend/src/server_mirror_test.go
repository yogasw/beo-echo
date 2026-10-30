package src

import (
	"testing"

	"beo-echo/backend/src/lib"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter_MirrorRegistersOnlyCatchAll(t *testing.T) {
	lib.MIRROR_PROJECT_ALIAS = "web-app"
	t.Cleanup(func() { lib.MIRROR_PROJECT_ALIAS = "" })

	routes := SetupRouter().Routes()
	assert.NotEmpty(t, routes)
	for _, r := range routes {
		assert.Equal(t, "/*path", r.Path, "%s %s must not exist on a mirror instance", r.Method, r.Path)
	}
}
