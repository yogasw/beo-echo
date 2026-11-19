package database

import "gorm.io/gorm"

// KeepDisabledOnCreate stores enabled=false on a record that was just created
// disabled. GORM leaves a zero-valued field that carries a default tag out of the
// INSERT, so a bool declared with gorm:"default:true" can never be created as
// false: the column default wins and the record goes live. Call it right after
// Create, in the same transaction, with the value the caller asked for.
func KeepDisabledOnCreate(tx *gorm.DB, model interface{}, enabled bool) error {
	if enabled {
		return nil
	}
	return tx.Model(model).Update("enabled", false).Error
}
