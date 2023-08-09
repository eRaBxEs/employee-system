// Package helper defines helper constants and functions used across other packages in the application
package helper

// SortOrder struct
type SortOrder string

// Key is a middleware key sting value
type Key string

const (
	// ZeroUUID default empty or non set UUID value
	ZeroUUID = "00000000-0000-0000-0000-000000000000"
	// LogStrKeyModule log service name value
	LogStrKeyModule = "ser_name"
	// LogStrKeyLevel log service level value
	LogStrKeyLevel = "lev_name"
	// LogStrKeyMethod log method name value
	LogStrKeyMethod = "method_name"
	// SortOrderASC for ascending sorting
	SortOrderASC = "ASC"
	// SortOrderDESC for descending sorting
	SortOrderDESC = "DESC"
)
