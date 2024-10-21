package errs

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorCode int

const (
	CodeInternal         ErrorCode = 1
	CodePermissionDenied ErrorCode = 2
	CodeNotFound         ErrorCode = 3
	CodeInvalidInput     ErrorCode = 4
)

func (ec ErrorCode) toGrpc() codes.Code {
	switch ec {
	case CodePermissionDenied:
		return codes.PermissionDenied
	case CodeNotFound:
		return codes.NotFound
	case CodeInvalidInput:
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

type CodableError struct {
	Code    ErrorCode
	Message string
}

func (e *CodableError) Error() string {
	return e.Message
}

func extractCodableErr(err error) *CodableError {
	for err != nil {
		ce, ok := err.(*CodableError)
		if ok {
			return ce
		}
		err = errors.Unwrap(err)
	}
	return &CodableError{
		Code: CodeInternal,
	}
}

func ToGrpcError(e error) error {
	msg := e.Error()
	ce := extractCodableErr(e)
	code := ce.Code.toGrpc()
	return status.Error(code, msg)
}
