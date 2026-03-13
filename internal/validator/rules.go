package validator

import (
	"fmt"
	"strings"
	"unicode"
)

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
			if r > unicode.MaxASCII {
				return fmt.Errorf("%w: found: %q", ErrOnlyEnglish, r)
			}
		}

		return nil
	}
}

func SpecialChars() Rule {
	return func(msg string) error {
		for _, r := range msg {
			if strings.ContainsRune("!@#$%^&*()+=[]{}|\\;<>?`~", r) {
				return fmt.Errorf("%w: found: %q", ErrSpecialChar, r)
			}
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
