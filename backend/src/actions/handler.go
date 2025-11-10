package actions

import (
	"beo-echo/backend/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ActionHandler handles HTTP requests for actions
type ActionHandler struct {
	service *ActionService
}

// NewActionHandler creates a new action handler
func NewActionHandler(service *ActionService) *ActionHandler {
	return &ActionHandler{service: service}
}

// CreateActionRequest represents the create action request
type CreateActionRequest struct {
	Name           string                      `json:"name" binding:"required"`
	Type           database.ActionType         `json:"type" binding:"required"`
	ExecutionPoint database.ExecutionPoint     `json:"execution_point"`
	Enabled        *bool                       `json:"enabled"` // Pointer to allow explicit false
	Priority       int                         `json:"priority"`
	Config         string                      `json:"config" binding:"required"`
	Filters        []CreateActionFilterRequest `json:"filters"`
}

// CreateActionFilterRequest represents the filter for an action
type CreateActionFilterRequest struct {
	Type     string `json:"type" binding:"required"`
	Key      string `json:"key"`
	Operator string `json:"operator" binding:"required"`
	Value    string `json:"value" binding:"required"`
}

// UpdateActionRequest represents the update action request
type UpdateActionRequest struct {
	Name           string                      `json:"name"`
	ExecutionPoint database.ExecutionPoint     `json:"execution_point"`
	Enabled        *bool                       `json:"enabled"`
	Priority       *int                        `json:"priority"`
	Config         string                      `json:"config"`
	Filters        []CreateActionFilterRequest `json:"filters"`
}

// CreateAction handles creating a new action
func (h *ActionHandler) CreateAction(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "project_id is required",
		})
		return
	}

	var req CreateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Set default execution point if not provided
	if req.ExecutionPoint == "" {
		req.ExecutionPoint = database.ExecutionPointAfterRequest
	}

	// Set default enabled if not provided
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	// Auto-set priority to last position (get max priority + 1)
	// Priority is 1-based, so first action gets priority 1
	priority := 1
	existingActions, err := h.service.GetActionsByProject(c.Request.Context(), projectID)
	if err == nil && len(existingActions) > 0 {
		// Find max priority
		maxPriority := 0
		for _, a := range existingActions {
			if a.Priority > maxPriority {
				maxPriority = a.Priority
			}
		}
		priority = maxPriority + 1
	}

	// Create action
	action := &database.Action{
		ProjectID:      projectID,
		Name:           req.Name,
		Type:           req.Type,
		ExecutionPoint: req.ExecutionPoint,
		Enabled:        enabled,
		Priority:       priority,
		Config:         req.Config,
	}

	if err := h.service.CreateAction(c.Request.Context(), action); err != nil {
		log.Error().Err(err).Msg("failed to create action")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Create filters if provided
	if len(req.Filters) > 0 {
		for _, filterReq := range req.Filters {
			filter := &database.ActionFilter{
				ActionID: action.ID,
				Type:     filterReq.Type,
				Key:      filterReq.Key,
				Operator: filterReq.Operator,
				Value:    filterReq.Value,
			}
			action.Filters = append(action.Filters, *filter)
		}
		// Update action with filters
		if err := h.service.UpdateAction(c.Request.Context(), action); err != nil {
			log.Error().Err(err).Msg("failed to create action filters")
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Action created successfully",
		"data":    action,
	})
}

// GetAction handles retrieving a single action
func (h *ActionHandler) GetAction(c *gin.Context) {
	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "action_id is required",
		})
		return
	}

	action, err := h.service.GetAction(c.Request.Context(), actionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Action not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    action,
	})
}

// GetProjectActions handles retrieving all actions for a project
func (h *ActionHandler) GetProjectActions(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "project_id is required",
		})
		return
	}

	actions, err := h.service.GetActionsByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve actions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    actions,
	})
}

// UpdateAction handles updating an existing action
func (h *ActionHandler) UpdateAction(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "action_id is required",
		})
		return
	}

	var req UpdateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Get existing action
	action, err := h.service.GetAction(c.Request.Context(), actionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Action not found",
		})
		return
	}

	// Update fields
	if req.Name != "" {
		action.Name = req.Name
	}
	if req.ExecutionPoint != "" {
		action.ExecutionPoint = req.ExecutionPoint
	}
	if req.Enabled != nil {
		action.Enabled = *req.Enabled
	}
	if req.Priority != nil {
		action.Priority = *req.Priority
	}
	if req.Config != "" {
		action.Config = req.Config
	}

	// Update filters if provided
	if req.Filters != nil {
		// Clear existing filters
		action.Filters = []database.ActionFilter{}

		// Add new filters
		for _, filterReq := range req.Filters {
			filter := database.ActionFilter{
				ActionID: action.ID,
				Type:     filterReq.Type,
				Key:      filterReq.Key,
				Operator: filterReq.Operator,
				Value:    filterReq.Value,
			}
			action.Filters = append(action.Filters, filter)
		}
	}

	if err := h.service.UpdateAction(c.Request.Context(), action); err != nil {
		log.Error().Err(err).Msg("failed to update action")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Action updated successfully",
		"data":    action,
	})
}

// DeleteAction handles deleting an action
func (h *ActionHandler) DeleteAction(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "action_id is required",
		})
		return
	}

	if err := h.service.DeleteAction(c.Request.Context(), actionID); err != nil {
		log.Error().Err(err).Msg("failed to delete action")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete action: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Action deleted successfully",
	})
}

// ToggleAction handles enabling/disabling an action
func (h *ActionHandler) ToggleAction(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "action_id is required",
		})
		return
	}

	// Get existing action
	action, err := h.service.GetAction(c.Request.Context(), actionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Action not found",
		})
		return
	}

	// Toggle enabled status
	action.Enabled = !action.Enabled

	if err := h.service.UpdateAction(c.Request.Context(), action); err != nil {
		log.Error().Err(err).Msg("failed to toggle action")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to toggle action: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Action toggled successfully",
		"data":    action,
	})
}

// UpdateActionPriorityRequest represents the request to update action priority
// Priority is 1-based (starts from 1, not 0)
type UpdateActionPriorityRequest struct {
	Priority int `json:"priority" binding:"required,min=1"`
}

// UpdateActionPriority handles updating the priority of an action
// This will automatically reorder other actions to maintain consistent priority ordering
func (h *ActionHandler) UpdateActionPriority(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "action_id is required",
		})
		return
	}

	var req UpdateActionPriorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Update priority via service (which handles reordering)
	if err := h.service.UpdateActionPriority(c.Request.Context(), actionID, req.Priority); err != nil {
		log.Error().Err(err).Msg("failed to update action priority")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Action priority updated successfully",
	})
}

// ActionTypeInfo represents information about an action type
type ActionTypeInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Category    string   `json:"category"`
	Fields      []string `json:"fields"`
}

// GetActionTypes returns all available action types
func (h *ActionHandler) GetActionTypes(c *gin.Context) {
	actionTypes := []ActionTypeInfo{
		{
			ID:          "replace_text",
			Name:        "Replace Text",
			Description: "Find and replace text in requests or responses using string matching or regex patterns",
			Icon:        "fa-exchange-alt",
			Category:    "Transform",
			Fields:      []string{"target", "pattern", "replacement", "use_regex", "header_key"},
		},
		{
			ID:          "run_javascript",
			Name:        "Run JavaScript",
			Description: "Execute custom JavaScript code to modify request or response data",
			Icon:        "fa-code",
			Category:    "Transform",
			Fields:      []string{"script"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    actionTypes,
	})
}
