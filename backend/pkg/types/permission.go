/*
Package types offers DTOs that cross layers.
*/
package types

import "errors"

const (
	// PermissionReadWrite used in file creation.
	PermissionReadWrite = 0600
)

var (
	// ErrCancelledContext for when the context was cancelled.
	ErrCancelledContext = errors.New("cancelled context")
)
