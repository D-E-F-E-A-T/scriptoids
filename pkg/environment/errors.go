package environment

import "errors"

var (
	EmptyIdentifiersError  = errors.New("package is missing a name and/or entry point")
	InvalidEntryPointError = errors.New("package has an invalid entry point")
	InvalidStateError      = errors.New("package is invalid, but can't determine why -- this is a bug")

	AlreadyLinkedError = errors.New("package is already linked")
	NotLinkedError     = errors.New("package is not linked")
)
