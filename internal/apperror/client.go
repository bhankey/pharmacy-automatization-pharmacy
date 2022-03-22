package apperror

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientErrorType int

const (
	Common ClientErrorType = iota
	WrongAuthorization
	WrongRequest
	WrongAuthToken
	NoClient
	WrongOneTimeCode
)

// thread safe, cus nobody writes.
var errorsMap = map[ClientErrorType]ClientError{
	Common:             errSomethingWentWrong,
	WrongRequest:       errWrongRequest,
	WrongAuthorization: errWrongAuthorization,
	WrongAuthToken:     errWrongAuthToken,
	NoClient:           errNoClient,
	WrongOneTimeCode:   errWrongOneTimeCode,
}

var errWrongAuthorization = ClientError{
	Code:    InvalidAuthorization,
	Message: "wrong password or email",
}

var errSomethingWentWrong = ClientError{
	Code:    Internal,
	Message: "something went wrong",
}

var errWrongRequest = ClientError{
	Code:    InvalidParams,
	Message: "wrong request",
}

var errWrongAuthToken = ClientError{
	Code:    InvalidAuthorization,
	Message: "wrong auth token",
}

var errNoClient = ClientError{
	Code:    InvalidParams,
	Message: "client doesn't exist",
}

var errWrongOneTimeCode = ClientError{
	Code:    InvalidParams,
	Message: "wrong one-time code",
}

type Code int

const (
	InvalidAuthorization Code = iota + 1
	InvalidParams
	PermissionDenied
	NotFound
	AlreadyExist
	Canceled
	Timeout
	Unavailable
	Internal
)

func (err ClientError) GetGRPCError() error {
	grpcErr := status.New(err.GetGRPCode(), err.Message)
	// grpcErr.WithDetails() TODO think about details https://cloud.google.com/apis/design/errors#error_model

	return grpcErr.Err() // nolint: wrapcheck, nolintlint
}

func (err ClientError) GetGRPCode() codes.Code {
	return map[Code]codes.Code{
		InvalidAuthorization: codes.Unauthenticated,
		InvalidParams:        codes.InvalidArgument,
		PermissionDenied:     codes.PermissionDenied,
		NotFound:             codes.NotFound,
		AlreadyExist:         codes.AlreadyExists,
		Canceled:             codes.Canceled,
		Timeout:              codes.DeadlineExceeded,
		Unavailable:          codes.Unavailable,
		Internal:             codes.Internal,
	}[err.Code]
}

func getCodeFromGRP(grpcCode codes.Code) Code {
	return map[codes.Code]Code{
		codes.Unauthenticated:  InvalidAuthorization,
		codes.InvalidArgument:  InvalidParams,
		codes.PermissionDenied: PermissionDenied,
		codes.NotFound:         NotFound,
		codes.AlreadyExists:    AlreadyExist,
		codes.Canceled:         Canceled,
		codes.DeadlineExceeded: Timeout,
		codes.Unavailable:      Unavailable,
		codes.Internal:         Internal,
	}[grpcCode]
}

func (err ClientError) GetHTTPCode() int {
	return map[Code]int{
		InvalidAuthorization: http.StatusUnauthorized,
		InvalidParams:        http.StatusBadRequest,
		PermissionDenied:     http.StatusForbidden,
		NotFound:             http.StatusNotFound,
		AlreadyExist:         http.StatusConflict,
		Canceled:             http.StatusRequestTimeout,
		Timeout:              http.StatusRequestTimeout,
		Unavailable:          http.StatusServiceUnavailable,
		Internal:             http.StatusInternalServerError,
	}[err.Code]
}

type ClientError struct {
	Code       Code
	Message    string
	ErrorToLog error
}

func NewClientError(errorType ClientErrorType, err error) ClientError {
	clientError, ok := errorsMap[errorType]
	if !ok {
		clientError = errorsMap[Common]
	}

	clientError.ErrorToLog = err

	return clientError
}

func NewClientErrorFromGRPC(err *status.Status) ClientError {
	return ClientError{
		Code:    getCodeFromGRP(err.Code()),
		Message: err.Message(),
	}
}

func (err ClientError) Error() string {
	return err.Message
}

func (err ClientError) Unwrap() error {
	return err.ErrorToLog
}
