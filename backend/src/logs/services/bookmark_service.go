package services

import (
	"beo-echo/backend/src/auth"
	"beo-echo/backend/src/database"
	systemConfig "beo-echo/backend/src/systemConfigs"
	"beo-echo/backend/src/utils"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// BookmarkService handles operations related to bookmarking logs
type BookmarkService struct {
	db *gorm.DB
}

// NewBookmarkService creates a new BookmarkService
func NewBookmarkService(db *gorm.DB) *BookmarkService {
	return &BookmarkService{
		db: db,
	}
}

// GetBookmarks retrieves all bookmarked logs for a project
func (s *BookmarkService) GetBookmarks(projectID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	if err := s.db.Table("request_logs").
		Where("project_id = ? AND bookmark = ?", projectID, true).
		Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// AddBookmark adds a log to bookmarks
func (s *BookmarkService) AddBookmark(projectID string, logData string) error {
	var logInput database.RequestLog
	if err := json.Unmarshal([]byte(logData), &logInput); err != nil {
		return errors.New("invalid log data: not a valid JSON")
	}

	// verification response hash
	if logInput.LogsHash == "" {
		return errors.New("invalid log data: response hash is empty")
	}

	logHash := logInput.LogsHash

	// verification jwt
	claim, errV := auth.ValidateToken(logInput.LogsHash)
	if errV != nil {
		log.Error().Err(errV).Msg("JWT validation failed")
		return errors.New("invalid log data: response hash is not a valid JWT")
	}

	// generate hash from input
	// remove log hash in input
	logInput.LogsHash = ""
	entry, err := json.Marshal(logInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal log entry to JSON")
	} else {
		md5Hash := utils.HashMD5(string(entry))
		if md5Hash != claim.UserID {
			return errors.New("invalid log data: response hash does not match")
		}
		logInput.LogsHash = logHash
	}

	var log database.RequestLog
	result := s.db.Where("project_id = ? AND id = ?", projectID, logInput.ID).First(&log)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logInput.Bookmark = true
			if err := s.db.Create(&logInput).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	} else {
		// Check if the log is already bookmarked
		if log.Bookmark {
			return errors.New("log is already bookmarked")
		} else {
			// Update the bookmark status
			logInput.Bookmark = true
			if err := s.db.Model(&log).Updates(logInput).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteBookmark removes a log from bookmarks
// when auto save is enabled update only bookmark field
// when auto save is disabled delete the log

func (s *BookmarkService) DeleteBookmark(projectID string, logID string) error {
	//check auto save is enabled or not
	autoSave, err := systemConfig.GetSystemConfigWithType[bool](systemConfig.AUTO_SAVE_LOGS_IN_DB_ENABLED)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get system config")
		return err
	}
	if autoSave {
		var log database.RequestLog
		result := s.db.Where("project_id = ? AND id = ?", projectID, logID).First(&log)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("log not found")
			}
			return result.Error
		}

		if !log.Bookmark {
			return errors.New("log is not bookmarked")
		}

		log.Bookmark = false
		if err := s.db.Save(&log).Error; err != nil {
			return err
		}
		return nil
	} else {
		result := s.db.Where("project_id = ? AND id = ?", projectID, logID).Delete(&database.RequestLog{})
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("log not found")
			}
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("log not found")
		}
		return nil
	}
}
