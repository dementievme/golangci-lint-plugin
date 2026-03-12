package validator

import "errors"

var (
	ErrLowerCase     = errors.New("log message should start with a lowercase letter")
	ErrOnlyEnglish   = errors.New("log message should be in English only")
	ErrSpecialChar   = errors.New("log message should not contain special character")
	ErrSensitiveData = errors.New("log message may expose sensitive data")
)
