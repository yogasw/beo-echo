package database

import "gorm.io/gorm"

// KeepDisabledOnCreate stores enabled=false on a record that was just created
// disabled. On Create, GORM replaces a zero-valued field that carries a default
// tag with that default, in both the INSERT and the struct, so a bool declared
// with gorm:"default:true" can never be created as false: the record goes live
// and the struct reads true. Call it right after Create, in the same transaction,
// with the value the caller asked for, then set the struct field back to it.
func KeepDisabledOnCreate(tx *gorm.DB, model interface{}, enabled bool) error {
	if enabled {
		return nil
	}
	return tx.Model(model).Update("enabled", false).Error
}
