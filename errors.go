package wsfanout

import "example.com/wsfanout/internal/errs"

var (
	ErrClosed       = errs.ErrClosed
	ErrNotFound     = errs.ErrNotFound
	ErrRoomFull     = errs.ErrRoomFull
	ErrAlreadyIn    = errs.ErrAlreadyIn
	ErrWriteFailed  = errs.ErrWriteFailed
	ErrEncodeFailed = errs.ErrEncodeFailed
	ErrNoEncoder    = errs.ErrNoEncoder
	ErrInvalidRoom  = errs.ErrInvalidRoom
	ErrInvalidConn  = errs.ErrInvalidConn
	ErrPersist      = errs.ErrPersist
	ErrCanceled     = errs.ErrCanceled
	ErrQueueFull    = errs.ErrQueueFull
)
