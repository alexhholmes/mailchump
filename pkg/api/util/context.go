package util

import "C"
import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

const (
	ContextUser   = "user"
	ContextLogger = "logger"
)

// Key is the type for all context.Context keys.
type Key string

func (k Key) String() string {
	return string(k)
}

func GetUserString(ctx context.Context) string {
	user, ok := ctx.Value(ContextUser).(string)
	if !ok {
		return uuid.Nil.String()
	}
	return user
}

func GetLogger(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(ContextLogger).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return log
}
