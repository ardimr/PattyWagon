package utils

import "strconv"

func ConvertIDString2Int(input string, defaultValue int64) int64 {
	if parsed, err := strconv.ParseInt(input, 10, 64); err == nil {
		return parsed
	}
	return defaultValue
}
