package util

const (
	ContextUser   = "user"
	ContextLogger = "logger"
)

// Key is the type for all context.Context keys.
type Key string

func (k Key) String() string {
	return string(k)
}
