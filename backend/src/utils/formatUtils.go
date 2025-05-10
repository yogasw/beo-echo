package utils

import (
	"fmt"
	"math"
)

// FormatFileSize converts a file size in bytes to a human-readable string (e.g., "1.2 KB")
func FormatFileSize(size int64) string {
	if size == 0 {
		return "0 B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	base := 1024.0
	i := math.Floor(math.Log(float64(size)) / math.Log(base))

	// Limit to the available units
	if int(i) >= len(units) {
		i = float64(len(units) - 1)
	}

	value := float64(size) / math.Pow(base, i)

	// Format with one decimal place for values less than 10
	if value < 10 {
		return fmt.Sprintf("%.1f %s", value, units[int(i)])
	}

	// Format without decimal places for larger values
	return fmt.Sprintf("%.0f %s", value, units[int(i)])
}
