package app

import (
	"errors"
	"testing"
)

// TestCodedError_Error_WithMessage tests Error() with message only
func TestCodedError_Error_WithMessage(t *testing.T) {
	err := &CodedError{
		Code:    "test_code",
		Message: "test message",
	}

	if err.Error() != "test message" {
		t.Fatalf("expected 'test message', got %q", err.Error())
	}
}

// TestCodedError_Error_WithWrappedError tests Error() with wrapped error
func TestCodedError_Error_WithWrappedError(t *testing.T) {
	inner := errors.New("inner error")
	err := &CodedError{
		Code: "test_code",
		Err:  inner,
	}

	if err.Error() != "inner error" {
		t.Fatalf("expected 'inner error', got %q", err.Error())
	}
}

// TestCodedError_Error_WithBoth tests Error() with both message and error
func TestCodedError_Error_WithBoth(t *testing.T) {
	inner := errors.New("inner")
	err := &CodedError{
		Code:    "test_code",
		Message: "outer",
		Err:     inner,
	}

	expected := "outer: inner"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

// TestCodedError_Error_Empty tests Error() with nil CodedError
func TestCodedError_Error_Empty(t *testing.T) {
	var err *CodedError = nil

	if err.Error() != "" {
		t.Fatalf("expected empty string for nil error, got %q", err.Error())
	}
}

// TestCodedError_Error_EmptyCodedError tests Error() with empty CodedError
func TestCodedError_Error_EmptyCodedError(t *testing.T) {
	err := &CodedError{}

	if err.Error() != "" {
		t.Fatalf("expected empty string, got %q", err.Error())
	}
}

// TestCodedError_Unwrap tests Unwrap() method
func TestCodedError_Unwrap(t *testing.T) {
	inner := errors.New("wrapped error")
	err := &CodedError{
		Code: "test_code",
		Err:  inner,
	}

	if unwrapped := err.Unwrap(); unwrapped != inner {
		t.Fatalf("expected wrapped error to be unwrapped, got %v", unwrapped)
	}
}

// TestCodedError_Unwrap_Nil tests Unwrap() with no wrapped error
func TestCodedError_Unwrap_Nil(t *testing.T) {
	err := &CodedError{
		Code: "test_code",
	}

	if unwrapped := err.Unwrap(); unwrapped != nil {
		t.Fatalf("expected nil from Unwrap(), got %v", unwrapped)
	}
}

// TestCodedError_Is tests errors.Is() support for CodedError
func TestCodedError_Is(t *testing.T) {
	inner := errors.New("test error")
	err := &CodedError{
		Code: "my_code",
		Err:  inner,
	}

	// errors.Is should work through Unwrap()
	if !errors.Is(err, inner) {
		t.Fatalf("expected errors.Is() to find wrapped error")
	}
}

// TestCodedError_As tests errors.As() support
func TestCodedError_As(t *testing.T) {
	inner := errors.New("test error")
	err := &CodedError{
		Code: "my_code",
		Err:  inner,
	}

	var target *CodedError
	if !errors.As(err, &target) {
		t.Fatalf("expected errors.As() to work for CodedError")
	}
	if target.Code != "my_code" {
		t.Fatalf("expected code 'my_code', got %q", target.Code)
	}
}

// TestCodedError_FieldAccess tests field access
func TestCodedError_FieldAccess(t *testing.T) {
	err := &CodedError{
		Code:    "error_code",
		Message: "error message",
		Err:     errors.New("inner"),
	}

	if err.Code != "error_code" {
		t.Fatalf("expected Code='error_code', got %q", err.Code)
	}
	if err.Message != "error message" {
		t.Fatalf("expected Message='error message', got %q", err.Message)
	}
	if err.Err == nil {
		t.Fatalf("expected Err to be non-nil")
	}
}

// TestCodedError_MessageOnly tests message-only error
func TestCodedError_MessageOnly(t *testing.T) {
	err := &CodedError{
		Code:    "validation_failed",
		Message: "invalid input format",
	}

	if err.Error() != "invalid input format" {
		t.Fatalf("expected 'invalid input format', got %q", err.Error())
	}
}

// TestCodedError_ChainedUnwrap tests multiple unwrap levels
func TestCodedError_ChainedUnwrap(t *testing.T) {
	innermost := errors.New("root cause")
	middle := &CodedError{
		Code: "mid_code",
		Err:  innermost,
	}
	outer := &CodedError{
		Code: "outer_code",
		Err:  middle,
	}

	// Should be able to find the innermost error through errors.Is
	if !errors.Is(outer, innermost) {
		t.Fatalf("expected to find innermost error through chain")
	}
}

// TestErrNotFound tests that ErrNotFound is properly set
func TestErrNotFound(t *testing.T) {
	if ErrNotFound == nil {
		t.Fatalf("expected ErrNotFound to be non-nil")
	}
}

// TestErrConflict tests that ErrConflict is properly set
func TestErrConflict(t *testing.T) {
	if ErrConflict == nil {
		t.Fatalf("expected ErrConflict to be non-nil")
	}
}

// TestCodedError_MessageWithSpecialChars tests message with special characters
func TestCodedError_MessageWithSpecialChars(t *testing.T) {
	err := &CodedError{
		Code:    "special_error",
		Message: "Error: \"something\" went wrong (code: 123)",
	}

	if err.Error() != "Error: \"something\" went wrong (code: 123)" {
		t.Fatalf("expected message with special chars, got %q", err.Error())
	}
}
