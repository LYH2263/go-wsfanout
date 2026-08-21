package payload

import "example.com/wsfanout/internal/errs"

// Validate 检查消息基本合法性。
func Validate(m Message) error {
	if m.Type == "" && len(m.Body) == 0 {
		return errs.ErrInvalidRoom
	}
	if len(m.Body) > 1<<20 {
		return errs.ErrEncodeFailed
	}
	return nil
}
