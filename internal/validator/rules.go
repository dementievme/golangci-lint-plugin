package validator

import (
	"fmt"
	"strings"
	"unicode"
)

const allowedPunctuation = ".,;:!?-–—'\"/()"

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

func isLatinLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

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
