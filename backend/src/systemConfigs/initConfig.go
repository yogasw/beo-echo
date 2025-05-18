package systemConfig

import (
	"beo-echo/backend/src/lib"
	"beo-echo/backend/src/utils"

	"github.com/rs/zerolog/log"
)

func InitializeDefaultConfig() error {
	key, err := GetSystemConfigWithType[string](JWT_SECRET)
	if err != nil {
		// Handle error
		log.Err(err).Msg("Error getting system config: " + err.Error())
	} else {
		if key == "" {
			key, _ = utils.GenerateRandomString(32)
			// Set default JWT secret
			err = SetSystemConfig(JWT_SECRET, key)
			if err != nil {
				// Handle error
				log.Err(err).Msg("Error setting system config: " + err.Error())
			}
		}
	}

	// init jwt secret
	lib.SetJWTSecret(key)
	return nil
}
