package utils

import (
	"database/sql"
	"strconv"
)

func ToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ToInt(ni sql.NullString) int64 {
	if ni.Valid {
		if i, err := strconv.ParseInt(ni.String, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func String2Int64(input string, defaultValue int64) int64 {
	if parsed, err := strconv.ParseInt(input, 10, 64); err == nil {
		return parsed
	}
	return defaultValue
}

func String2Int32(input string, defaultValue int) int {
	if parsed, err := strconv.ParseInt(input, 10, 32); err == nil {
		return int(parsed)
	}
	return defaultValue
}
