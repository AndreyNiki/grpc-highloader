package ptr

// ToPtr return pointer to comparable type.
func ToPtr[T comparable](value T) *T {
	return &value
}
