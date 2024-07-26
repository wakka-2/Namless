package repository

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

// NoopLogger offers a logger that does nothing.
//
// Used in tests, so that repo operations do not clog the console with warnings (slow operation took longer than 200ms).
type NoopLogger struct{}

// NewNoopLogger builds a new noop logger.
func NewNoopLogger() logger.Interface {
	return &NoopLogger{}
}

// LogMode builds a new noop logger.
func (l *NoopLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return &NoopLogger{}
}

// Info does nothing.
func (l *NoopLogger) Info(context.Context, string, ...interface{}) {}

// Warn does nothing.
func (l *NoopLogger) Warn(context.Context, string, ...interface{}) {}

// Error does nothing.
func (l *NoopLogger) Error(context.Context, string, ...interface{}) {}

// Trace does nothing.
func (l *NoopLogger) Trace(_ context.Context, _ time.Time, _ func() (sql string, rowsAffected int64), _ error) {
}
