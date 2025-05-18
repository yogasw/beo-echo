package handler

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/logs/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var bookmarkService *services.BookmarkService

// InitBookmarkService initializes the bookmark service
func InitBookmarkService() {
	db := database.GetDB()
	if db == nil {
		return
	}

	bookmarkService = services.NewBookmarkService(db)
}

// EnsureBookmarkService ensures bookmark service is initialized
func EnsureBookmarkService() {
	if bookmarkService == nil {
		InitBookmarkService()
	}
}

// AddBookmarkHandler handles saving a log as bookmark
func AddBookmarkHandler(c *gin.Context) {
	EnsureBookmarkService()
	if bookmarkService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Bookmark service is not available",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Parse request body
	var request struct {
		LogStr string `json:"logs" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Save bookmark
	err := bookmarkService.AddBookmark(projectID, request.LogStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Error saving bookmark: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bookmark added successfully",
	})
}

// DeleteBookmarkHandler handles deleting a bookmark
func DeleteBookmarkHandler(c *gin.Context) {
	EnsureBookmarkService()
	if bookmarkService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Bookmark service is not available",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	bookmarkID := c.Param("bookmarkId")
	if bookmarkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Bookmark id is required",
		})
		return
	}

	// Delete bookmark
	if err := bookmarkService.DeleteBookmark(projectID, bookmarkID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error deleting bookmark: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bookmark deleted successfully",
	})
}

// GetBookmarksHandler handles retrieving all bookmarks
func GetBookmarksHandler(c *gin.Context) {
	EnsureBookmarkService()
	if bookmarkService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Bookmark service is not available",
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Get bookmarks
	bookmarks, err := bookmarkService.GetBookmarks(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error retrieving bookmarks: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookmarks,
	})
}
