package error

import (
	"github.com/pkg/errors"
)

// Errors returned on rsmq operation
var (
	ErrNoAttributeSupplied  = errors.New("no attribute supplied")
	ErrMissingParameter     = errors.New("missing parameter")
	ErrInvalidFormat        = errors.New("invalid format")
	ErrInvalidValue         = errors.New("invalid value")
	ErrQueueExists          = errors.New("queue already exists")
	ErrQueueEmpty           = errors.New("queue is empty")
	ErrQueueNotFound        = errors.New("queue not found")
	ErrQueueFull			= errors.New("queue is full")
	ErrMessageNotFound      = errors.New("message not found")
	ErrMessageExists        = errors.New("message already exists")
	ErrMessageTooLong       = errors.New("message too long")
	ErrMessageNotString     = errors.New("message is not a string")
)