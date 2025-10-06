package errors

import (
	"fmt"
	"strings"
)

// ErrorCode represents specific error types
type ErrorCode string

const (
	// Documentation errors
	ErrDocNotFound      ErrorCode = "DOC_NOT_FOUND"
	ErrDocInvalid       ErrorCode = "DOC_INVALID"
	ErrDocVersionNotFound ErrorCode = "DOC_VERSION_NOT_FOUND"
	
	// Package errors
	ErrPackageNotFound  ErrorCode = "PACKAGE_NOT_FOUND"
	ErrCategoryNotFound ErrorCode = "CATEGORY_NOT_FOUND"
	ErrCatalogInvalid   ErrorCode = "CATALOG_INVALID"
	
	// Update errors
	ErrUpdateFailed     ErrorCode = "UPDATE_FAILED"
	ErrGitHubAPI        ErrorCode = "GITHUB_API_ERROR"
	
	// External errors
	ErrFetchFailed      ErrorCode = "FETCH_FAILED"
	ErrURLInvalid       ErrorCode = "URL_INVALID"
	ErrContentTooLarge  ErrorCode = "CONTENT_TOO_LARGE"
	
	// Generic errors
	ErrInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrInternal         ErrorCode = "INTERNAL_ERROR"
)

// MCPError represents a structured error for MCP tools
type MCPError struct {
	Code    ErrorCode
	Message string
	Details map[string]string
	Cause   error
}

// Error implements the error interface
func (e *MCPError) Error() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", e.Code))
	parts = append(parts, e.Message)
	
	if len(e.Details) > 0 {
		var details []string
		for k, v := range e.Details {
			details = append(details, fmt.Sprintf("%s=%s", k, v))
		}
		parts = append(parts, fmt.Sprintf("(%s)", strings.Join(details, ", ")))
	}
	
	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("- %v", e.Cause))
	}
	
	return strings.Join(parts, " ")
}

// Unwrap returns the underlying cause
func (e *MCPError) Unwrap() error {
	return e.Cause
}

// New creates a new MCPError
func New(code ErrorCode, message string) *MCPError {
	return &MCPError{
		Code:    code,
		Message: message,
		Details: make(map[string]string),
	}
}

// Wrap wraps an existing error with context
func Wrap(code ErrorCode, message string, cause error) *MCPError {
	return &MCPError{
		Code:    code,
		Message: message,
		Details: make(map[string]string),
		Cause:   cause,
	}
}

// WithDetail adds a detail to the error
func (e *MCPError) WithDetail(key, value string) *MCPError {
	e.Details[key] = value
	return e
}

// UserFriendlyMessage returns a user-friendly error message
func UserFriendlyMessage(err error) string {
	if mcpErr, ok := err.(*MCPError); ok {
		switch mcpErr.Code {
		case ErrDocNotFound:
			return "Documentation file not found. Please check the filename and version."
		case ErrDocVersionNotFound:
			return "The specified Laravel version is not available. Supported versions: 12.x, 11.x, 10.x, 9.x, 8.x, 7.x, 6.x"
		case ErrPackageNotFound:
			return "Package not found in the catalog. Try searching for similar packages."
		case ErrCategoryNotFound:
			return "Category not found. Use list_package_categories to see available categories."
		case ErrURLInvalid:
			return "Invalid URL. Please provide a valid HTTP or HTTPS URL."
		case ErrContentTooLarge:
			return "The external resource is too large to fetch (max 5MB)."
		case ErrGitHubAPI:
			return "Failed to connect to GitHub API. Please check your internet connection or try again later."
		default:
			return mcpErr.Message
		}
	}
	return err.Error()
}
