package util

// Key is the type for all context.Context keys.
type Key string

func (k Key) String() string {
	return string(k)
}
