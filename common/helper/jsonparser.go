package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	MaxPayloadSize = 5 << 20 // 5 MB
	MaxArrayItems  = 1000
	MaxNameLength  = 255
)

var (
	ErrPayloadTooLarge  = errors.New("payload too large")
	ErrInvalidUTF8      = errors.New("invalid utf-8")
	ErrTrailingData     = errors.New("unexpected trailing data")
	ErrDuplicateField   = errors.New("duplicate field detected")
	ErrTooManyItems     = errors.New("too many items")
	ErrInvalidRoot      = errors.New("root must be json array")
	ErrEmptyName        = errors.New("name is required")
	ErrNameTooLong      = errors.New("name too long")
	ErrInvalidUUID      = errors.New("invalid uuid")
	ErrInvalidFieldName = errors.New("invalid field name")
)

// rejectDuplicateKeys:
// recursive duplicate detector using tokenizer
func RejectDuplicateKeys(raw []byte) error {
	dec := json.NewDecoder(bytes.NewReader(raw))

	type objectFrame struct {
		keys map[string]struct{}
	}

	var stack []objectFrame
	var expectingKey bool

	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("token parse error: %w", err)
		}

		switch t := token.(type) {

		case json.Delim:
			switch t {

			case '{':
				stack = append(stack, objectFrame{
					keys: make(map[string]struct{}),
				})
				expectingKey = true

			case '}':
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
				expectingKey = false

			case '[':
				expectingKey = false

			case ']':
				expectingKey = false
			}

		case string:
			if len(stack) == 0 {
				continue
			}

			if expectingKey {
				key := normalizeKey(t)

				if !isValidFieldName(key) {
					return fmt.Errorf("%w: %s", ErrInvalidFieldName, t)
				}

				top := &stack[len(stack)-1]

				if _, exists := top.keys[key]; exists {
					return fmt.Errorf("%w: %s", ErrDuplicateField, t)
				}

				top.keys[key] = struct{}{}
				expectingKey = false

			} else {
				expectingKey = true
			}

		default:
			if len(stack) > 0 {
				expectingKey = true
			}
		}
	}

	return nil
}

func normalizeKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func isValidFieldName(s string) bool {
	if s == "" {
		return false
	}

	// prevent weird pollution chars
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '_':
		default:
			return false
		}
	}

	return true
}
