// Copyright (c) 2026 Michael_Go. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package validator

import "errors"

// Sentinel errors returned by validation rules.
var (
	// ErrLowerCase indicates the message starts with an uppercase letter.
	ErrLowerCase = errors.New("log message should start with a lowercase letter")
	// ErrOnlyEnglish indicates the message contains non-Latin letters.
	ErrOnlyEnglish = errors.New("log message should be in English only")
	// ErrSpecialChar indicates the message contains a disallowed special character.
	ErrSpecialChar = errors.New("log message should not contain special character")
	// ErrSensitiveData indicates the message may expose sensitive information.
	ErrSensitiveData = errors.New("log message may expose sensitive data")
)
