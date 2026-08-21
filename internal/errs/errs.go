package errs

import (
	"errors"
	"fmt"
)

var (
	ErrClosed       = errors.New("wsfanout: hub closed")
	ErrNotFound     = errors.New("wsfanout: not found")
	ErrRoomFull     = errors.New("wsfanout: room full")
	ErrAlreadyIn    = errors.New("wsfanout: already joined")
	ErrWriteFailed  = errors.New("wsfanout: write failed")
	ErrEncodeFailed = errors.New("wsfanout: encode failed")
	ErrNoEncoder    = errors.New("wsfanout: encoder is nil")
	ErrInvalidRoom  = errors.New("wsfanout: invalid room")
	ErrInvalidConn  = errors.New("wsfanout: invalid conn")
	ErrPersist      = errors.New("wsfanout: persist failed")
	ErrCanceled     = errors.New("wsfanout: canceled")
	ErrQueueFull    = errors.New("wsfanout: queue full")
	ErrTimeout      = errors.New("wsfanout: timeout")
	ErrAudit        = errors.New("wsfanout: audit failed")
)

// WrapWrite 将底层写错误包装为 ErrWriteFailed（双 %w，哨兵与原因均可 errors.Is）。
func WrapWrite(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrWriteFailed, err)
}

// WrapEncode 包装编码错误。
func WrapEncode(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", ErrEncodeFailed, err)
}

// WrapPersist 包装持久化错误。
func WrapPersist(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", ErrPersist, err)
}

// WrapCanceled 包装取消。
func WrapCanceled(err error) error {
	if err == nil {
		return ErrCanceled
	}
	return fmt.Errorf("%w: %w", ErrCanceled, err)
}

// WrapAudit 包装审计错误。
func WrapAudit(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", ErrAudit, err)
}
