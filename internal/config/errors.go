package config

import "errors"

var (
	ErrConfigPathIsEmpty = errors.New("config path is empty")
	ErrPathNotSpecified  = errors.New("config path is not specified")
)
