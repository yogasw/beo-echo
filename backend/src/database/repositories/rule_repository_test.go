package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

func TestRuleRepository_RegexValue(t *testing.T) {
	database.SetupTestEnvironment(t)
	db := database.GetDB()

	project := &database.Project{ID: uuid.New().String(), Name: "Test Project"}
	require.NoError(t, db.Create(project).Error)
	endpoint := &database.MockEndpoint{ID: uuid.New().String(), ProjectID: project.ID, Path: "/test", Method: "GET"}
	require.NoError(t, db.Create(endpoint).Error)
	response := &database.MockResponse{ID: uuid.New().String(), EndpointID: endpoint.ID, StatusCode: 200}
	require.NoError(t, db.Create(response).Error)

	repo := NewRuleRepository(db)

	t.Run("create rejects a pattern that does not compile", func(t *testing.T) {
		err := repo.CreateRule(&database.MockRule{ResponseID: response.ID, Type: "query", Key: "q", Operator: "regex", Value: "(unclosed"})
		assert.ErrorContains(t, err, "invalid regex value")
	})

	t.Run("create accepts a valid pattern", func(t *testing.T) {
		rule := &database.MockRule{ResponseID: response.ID, Type: "query", Key: "q", Operator: "regex", Value: `^\d{8}$`}
		require.NoError(t, repo.CreateRule(rule))

		rule.Value = "[bad"
		assert.ErrorContains(t, repo.UpdateRule(rule), "invalid regex value")
	})

	t.Run("non-regex operators keep accepting any value", func(t *testing.T) {
		assert.NoError(t, repo.CreateRule(&database.MockRule{ResponseID: response.ID, Type: "query", Key: "q", Operator: "contains", Value: "(unclosed"}))
	})
}
