package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppErrorFormatting(t *testing.T) {
	err := New(ErrInvalidInput, "bad request")
	require.Equal(t, "[INVALID_INPUT] bad request", err.Error())

	errWithCause := NewWithErr(ErrInternal, "failed", errors.New("db down"))
	require.Equal(t, "[INTERNAL_ERROR] failed: db down", errWithCause.Error())
}
