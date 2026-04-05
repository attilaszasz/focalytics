package app

import (
	"errors"
	"fmt"
)

type ExitPolicy struct {
	Success        int
	InvalidInput   int
	RuntimeFailure int
}

type CommandError struct {
	Code int
	Err  error
}

func (e *CommandError) Error() string {
	return e.Err.Error()
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

func DefaultExitPolicy() ExitPolicy {
	return ExitPolicy{
		Success:        0,
		InvalidInput:   2,
		RuntimeFailure: 1,
	}
}

func CommandErrorFromCode(code int, err error) error {
	if err == nil {
		err = fmt.Errorf("command failed with exit code %d", code)
	}

	return &CommandError{Code: code, Err: err}
}

func InvalidInputError(err error) error {
	return CommandErrorFromCode(DefaultExitPolicy().InvalidInput, err)
}

func RuntimeFailureError(err error) error {
	return CommandErrorFromCode(DefaultExitPolicy().RuntimeFailure, err)
}

func ExitCodeForError(err error, policy ExitPolicy) int {
	if err == nil {
		return policy.Success
	}

	var coded *CommandError
	if errors.As(err, &coded) {
		return coded.Code
	}

	return policy.RuntimeFailure
}
