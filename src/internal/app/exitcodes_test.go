package app

import (
	"errors"
	"testing"
)

func TestExitCodeForErrorUsesCommandErrorCode(t *testing.T) {
	policy := DefaultExitPolicy()
	err := InvalidInputError(errors.New("bad input"))

	if got := ExitCodeForError(err, policy); got != policy.InvalidInput {
		t.Fatalf("unexpected exit code: got %d want %d", got, policy.InvalidInput)
	}
}

func TestExitCodeForErrorFallsBackToRuntimeFailure(t *testing.T) {
	policy := DefaultExitPolicy()

	if got := ExitCodeForError(errors.New("boom"), policy); got != policy.RuntimeFailure {
		t.Fatalf("unexpected exit code: got %d want %d", got, policy.RuntimeFailure)
	}
}

func TestRuntimeFailureErrorUsesRuntimeCode(t *testing.T) {
	policy := DefaultExitPolicy()
	err := RuntimeFailureError(errors.New("runtime failed"))

	if got := ExitCodeForError(err, policy); got != policy.RuntimeFailure {
		t.Fatalf("unexpected exit code: got %d want %d", got, policy.RuntimeFailure)
	}
}

func TestCommandErrorFromCodeUsesFallbackMessage(t *testing.T) {
	err := CommandErrorFromCode(9, nil)

	if err.Error() == "" {
		t.Fatal("expected fallback error message")
	}
}
