package utils

// ToPointer takes a value of any type and returns a pointer to that value
func ToPointer[T any](val T) *T {
	return &val
}

// PointerValue :
// this func will return default value or value of variable
// to handle panic if the value of variable is nil
func PointerValue[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}
