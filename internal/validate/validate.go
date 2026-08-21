package validate

import (
	"strings"
	"unicode"

	"example.com/wsfanout/internal/errs"
)

// RoomName 校验房间名。
func RoomName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 128 {
		return errs.ErrInvalidRoom
	}
	for _, r := range name {
		if unicode.IsSpace(r) {
			return errs.ErrInvalidRoom
		}
		if r == '/' || r == '\\' {
			return errs.ErrInvalidRoom
		}
	}
	return nil
}

// ConnID 校验连接 ID。
func ConnID(id string) error {
	if id == "" || len(id) > 128 {
		return errs.ErrInvalidConn
	}
	for _, r := range id {
		if unicode.IsSpace(r) {
			return errs.ErrInvalidConn
		}
	}
	return nil
}

// MessageType 校验类型名。
func MessageType(t string) error {
	if len(t) > 64 {
		return errs.ErrEncodeFailed
	}
	return nil
}
