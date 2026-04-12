package runtime

import "strings"

type ErrorCode string

const (
	ErrScopeInvalid      ErrorCode = "E_SCOPE_INVALID"
	ErrSymbolNotFound    ErrorCode = "E_SYMBOL_NOT_FOUND"
	ErrDuplicateSymbol   ErrorCode = "E_DUPLICATE_SYMBOL"
	ErrTypeInvalid       ErrorCode = "E_TYPE_INVALID"
	ErrChainInvalid      ErrorCode = "E_CHAIN_INVALID"
	ErrHandlerUnresolved ErrorCode = "E_HANDLER_UNRESOLVED"
)

type AGSError struct {
	Code         ErrorCode
	Message      string
	SemanticPath string
	Line         int
}

func (e *AGSError) Error() string {
	parts := []string{string(e.Code), e.Message}
	if strings.TrimSpace(e.SemanticPath) != "" {
		parts = append(parts, "path="+e.SemanticPath)
	}
	if e.Line > 0 {
		parts = append(parts, "line="+itoa(e.Line))
	}
	return strings.Join(parts, " | ")
}

type AGSMultiError struct {
	Errors []*AGSError
}

func (e *AGSMultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	parts := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		parts = append(parts, err.Error())
	}
	return strings.Join(parts, "\n")
}

func (e *AGSMultiError) Codes() []ErrorCode {
	codes := make([]ErrorCode, 0, len(e.Errors))
	for _, err := range e.Errors {
		codes = append(codes, err.Code)
	}
	return codes
}

func newAGSError(code ErrorCode, message string, semanticPath string) *AGSError {
	return &AGSError{
		Code:         code,
		Message:      message,
		SemanticPath: semanticPath,
		Line:         0,
	}
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	digits := []byte{}
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}

	return string(digits)
}
