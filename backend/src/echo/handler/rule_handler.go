package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/services"
)

// RuleHandler handles HTTP requests for rules
type RuleHandler struct {
	service *services.RuleService
}

// NewRuleHandler creates a new rule handler
func NewRuleHandler(ruleService *services.RuleService) *RuleHandler {
	return &RuleHandler{
		service: ruleService,
	}
}

// ListRulesHandler lists all rules for a response
//
// Sample curl:
// curl -X GET "http://localhost:3600/api/api/projects/my-project/endpoints/endpoint-id/responses/response-id/rules" -H "Content-Type: application/json"
func (h *RuleHandler) ListRulesHandler(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Rule service not initialized",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	responseID := c.Param("responseId")
	if responseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Response ID is required",
		})
		return
	}

	// Validate hierarchy using the service layer
	isValid, err := h.service.ValidateResponseHierarchy(projectID, endpointID, responseID)
	if err != nil || !isValid {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Invalid hierarchy: Response not found or doesn't belong to the specified endpoint/project",
		})
		return
	}

	// Get rules
	rules, err := h.service.GetRules(responseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve rules: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rules,
	})
}

// GetRuleHandler retrieves a rule by ID
//
// Sample curl:
// curl -X GET "http://localhost:3600/api/api/projects/my-project/endpoints/endpoint-id/responses/response-id/rules/rule-id" -H "Content-Type: application/json"
func (h *RuleHandler) GetRuleHandler(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Rule service not initialized",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	responseID := c.Param("responseId")
	if responseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Response ID is required",
		})
		return
	}

	ruleID := c.Param("ruleId")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Rule ID is required",
		})
		return
	}

	// Validate hierarchy using the service layer
	isValid, err := h.service.ValidateResponseHierarchy(projectID, endpointID, responseID)
	if err != nil || !isValid {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Invalid hierarchy: Response not found or doesn't belong to the specified endpoint/project",
		})
		return
	}

	// Get rule and validate it belongs to the response
	rule, err := h.service.GetRule(ruleID, responseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Rule not found or doesn't belong to the specified response",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rule,
	})
}

// CreateRuleHandler creates a new rule
//
// Sample curl:
//
//	curl -X POST "http://localhost:3600/api/api/projects/my-project/endpoints/endpoint-id/responses/response-id/rules" \
//	  -H "Content-Type: application/json" \
//	  -d '{
//	    "type": "header",
//	    "key": "X-Auth-Token",
//	    "operator": "equals",
//	    "value": "secret-token"
//	  }'
func (h *RuleHandler) CreateRuleHandler(c *gin.Context) {

	if h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Rule service not initialized",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	responseID := c.Param("responseId")
	if responseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Response ID is required",
		})
		return
	}

	// Validate hierarchy using the service layer
	isValid, err := h.service.ValidateResponseHierarchy(projectID, endpointID, responseID)
	if err != nil || !isValid {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Invalid hierarchy: Response not found or doesn't belong to the specified endpoint/project",
		})
		return
	}

	// Parse rule data from request body
	var rule database.MockRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid rule data: " + err.Error(),
		})
		return
	}

	// Set response ID from URL parameter
	rule.ResponseID = responseID

	// Create rule through service layer
	createdRule, err := h.service.CreateRule(&rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create rule: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Rule created successfully",
		"data":    createdRule,
	})
}

// UpdateRuleHandler updates an existing rule
//
// Sample curl:
//
//	curl -X PUT "http://localhost:3600/api/api/projects/my-project/endpoints/endpoint-id/responses/response-id/rules/rule-id" \
//	  -H "Content-Type: application/json" \
//	  -d '{
//	    "type": "query",
//	    "key": "token",
//	    "operator": "contains",
//	    "value": "new-value"
//	  }'
func (h *RuleHandler) UpdateRuleHandler(c *gin.Context) {

	if h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Rule service not initialized",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	responseID := c.Param("responseId")
	if responseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Response ID is required",
		})
		return
	}

	ruleID := c.Param("ruleId")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Rule ID is required",
		})
		return
	}

	// Validate hierarchy and rule ownership using the service layer
	isValid, err := h.service.ValidateRuleHierarchy(projectID, endpointID, responseID, ruleID)
	if err != nil || !isValid {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Invalid hierarchy: Rule not found or doesn't belong to the specified response/endpoint/project",
		})
		return
	}

	// Parse update data from request body
	var updateData database.MockRule
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid rule data: " + err.Error(),
		})
		return
	}

	// Update rule through service layer
	updatedRule, err := h.service.UpdateRule(ruleID, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update rule: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rule updated successfully",
		"data":    updatedRule,
	})
}

// DeleteRuleHandler deletes a rule
//
// Sample curl:
// curl -X DELETE "http://localhost:3600/api/api/projects/my-project/endpoints/endpoint-id/responses/response-id/rules/rule-id" -H "Content-Type: application/json"
func (h *RuleHandler) DeleteRuleHandler(c *gin.Context) {

	if h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Rule service not initialized",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	responseID := c.Param("responseId")
	if responseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Response ID is required",
		})
		return
	}

	ruleID := c.Param("ruleId")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Rule ID is required",
		})
		return
	}

	// Validate hierarchy and rule ownership using the service layer
	isValid, err := h.service.ValidateRuleHierarchy(projectID, endpointID, responseID, ruleID)
	if err != nil || !isValid {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Invalid hierarchy: Rule not found or doesn't belong to the specified response/endpoint/project",
		})
		return
	}

	// Delete rule through service layer
	err = h.service.DeleteRule(ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete rule: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rule deleted successfully",
	})
}

// DeleteAllRulesHandler deletes all rules for a response
//
// Sample curl:
// curl -X DELETE "http://localhost:3600/api/api/projects/my-project/endpoints/endpoint-id/responses/response-id/rules" -H "Content-Type: application/json"
func (h *RuleHandler) DeleteAllRulesHandler(c *gin.Context) {
	if h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Rule service not initialized",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	responseID := c.Param("responseId")
	if responseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Response ID is required",
		})
		return
	}

	// Validate hierarchy using the service layer
	isValid, err := h.service.ValidateResponseHierarchy(projectID, endpointID, responseID)
	if err != nil || !isValid {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Invalid hierarchy: Response not found or doesn't belong to the specified endpoint/project",
		})
		return
	}

	// Delete all rules through service layer
	err = h.service.DeleteRulesByResponse(responseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete rules: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All rules deleted successfully",
	})
}
