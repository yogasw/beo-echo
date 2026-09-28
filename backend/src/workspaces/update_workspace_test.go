package workspaces

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
)

// updateRepoStub records what UpdateWorkspace was called with. Every other
// method of WorkspaceRepository is left nil on the embedded interface, so a
// test that reaches one panics instead of quietly passing.
type updateRepoStub struct {
	WorkspaceRepository

	calls     int
	gotID     string
	gotName   string
	returnWS  *database.Workspace
	returnErr error
}

func (s *updateRepoStub) UpdateWorkspace(ctx context.Context, workspaceID string, name string) (*database.Workspace, error) {
	s.calls++
	s.gotID = workspaceID
	s.gotName = name
	return s.returnWS, s.returnErr
}

func TestUpdateWorkspace_TrimsAndSaves(t *testing.T) {
	// Given - a repository that accepts the rename
	stub := &updateRepoStub{returnWS: &database.Workspace{ID: "ws-1", Name: "test ss"}}
	svc := NewWorkspaceService(stub)

	// When - the name arrives with surrounding whitespace
	got, err := svc.UpdateWorkspace(context.Background(), "ws-1", "  test ss  ")

	// Then - it is trimmed before it reaches the repository
	require.NoError(t, err)
	assert.Equal(t, 1, stub.calls)
	assert.Equal(t, "ws-1", stub.gotID)
	assert.Equal(t, "test ss", stub.gotName)
	assert.Equal(t, "test ss", got.Name)
}

func TestUpdateWorkspace_RejectsBlankName(t *testing.T) {
	for _, name := range []string{"", "   ", "\t\n"} {
		// Given - a repository that must not be touched
		stub := &updateRepoStub{}
		svc := NewWorkspaceService(stub)

		// When - a blank name is submitted
		_, err := svc.UpdateWorkspace(context.Background(), "ws-1", name)

		// Then - it is rejected before any write
		assert.ErrorIs(t, err, ErrWorkspaceNameRequired, "name %q", name)
		assert.Zero(t, stub.calls, "repository must not be called for name %q", name)
	}
}

func TestUpdateWorkspace_RejectsOverlongName(t *testing.T) {
	// Given - a name one rune past the limit
	stub := &updateRepoStub{}
	svc := NewWorkspaceService(stub)

	// When
	_, err := svc.UpdateWorkspace(context.Background(), "ws-1", strings.Repeat("a", maxWorkspaceNameLength+1))

	// Then
	assert.ErrorIs(t, err, ErrWorkspaceNameTooLong)
	assert.Zero(t, stub.calls)
}

func TestUpdateWorkspace_AcceptsNameAtTheLimit(t *testing.T) {
	// Given - a name exactly at the limit, counted in runes not bytes
	stub := &updateRepoStub{returnWS: &database.Workspace{ID: "ws-1"}}
	svc := NewWorkspaceService(stub)

	// When - the name is multi-byte throughout
	_, err := svc.UpdateWorkspace(context.Background(), "ws-1", strings.Repeat("é", maxWorkspaceNameLength))

	// Then - byte length is irrelevant, it is accepted
	require.NoError(t, err)
	assert.Equal(t, 1, stub.calls)
}

func TestUpdateWorkspace_PropagatesNotFound(t *testing.T) {
	// Given - a repository that cannot find the workspace
	stub := &updateRepoStub{returnErr: gorm.ErrRecordNotFound}
	svc := NewWorkspaceService(stub)

	// When
	_, err := svc.UpdateWorkspace(context.Background(), "missing", "whatever")

	// Then - the handler can still tell a 404 apart from a 500
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
