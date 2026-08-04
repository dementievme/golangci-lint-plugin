// Copyright (c) 2026 Michael_Go. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package validator

import (
	"fmt"
	"strings"
	"unicode"
)

// allowedPunctuation defines the set of punctuation characters
// permitted in log messages by the SpecialChars rule.
const allowedPunctuation = ".,;:!?-–—'\"/()"

// Lowercase returns a rule that checks whether the log message
// starts with a lowercase letter.
func Lowercase() Rule {
	return func(msg string) error {
		if msg == "" {
			return nil
		}

		if unicode.IsUpper([]rune(msg)[0]) {
			return ErrLowerCase
		}

		return nil
	}
}

// English returns a rule that checks whether all letter characters
// in the message are Latin (a-z, A-Z). Non-letter characters
// (digits, spaces, punctuation) are ignored.
func English() Rule {
	return func(msg string) error {
		for _, r := range msg {
			if unicode.IsLetter(r) && !isLatinLetter(r) {
				return fmt.Errorf("%w: found: %q", ErrOnlyEnglish, r)
			}
		}

		return nil
	}
}

// isLatinLetter reports whether r is a basic Latin letter (a-z or A-Z).
func isLatinLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// SpecialChars returns a rule that rejects characters that are not
// letters, digits, whitespace, or allowed punctuation.
// This catches emoji, programming symbols (@#$%^&*), and other
// non-standard characters.
func SpecialChars() Rule {
	return func(msg string) error {
		for _, r := range msg {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
				continue
			}
			if strings.ContainsRune(allowedPunctuation, r) {
				continue
			}
			return fmt.Errorf("%w: found: %q", ErrSpecialChar, r)
		}

		return nil
	}
}

// SensitiveData returns a rule that checks whether the message
// contains any of the given sensitive keywords (case-insensitive).
func SensitiveData(keywords []string) Rule {
	return func(msg string) error {
		lower := strings.ToLower(msg)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				return fmt.Errorf("%w: found: %q", ErrSensitiveData, kw)
			}
		}

		return nil
	}
}
