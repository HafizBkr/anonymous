package commons

import (
	"anonymous/types"
	"errors"
	"net/http"
)

type errs struct {
	ResourceNotFound     error
	InternalServerError  error
	Conflict             error
	TokenEncodingFailed  error
	DuplicateEmail       error
	AuthenticationFailed error
	DuplicateUsername    error
}

var Errors = errs{
	ResourceNotFound: errors.New("ResourceNotFound"),
	InternalServerError: types.ServiceError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  Codes.InternalError,
	},
	Conflict: types.ServiceError{
		StatusCode: http.StatusConflict,
		ErrorCode:  "Conflict",
	},
	TokenEncodingFailed: types.ServiceError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "TokenEncodingFailed",
	},
	DuplicateEmail: types.ServiceError{
		StatusCode: http.StatusConflict,
		ErrorCode:  Codes.DuplicateEmail,
	},
	DuplicateUsername: types.ServiceError{
			StatusCode: http.StatusConflict,
			ErrorCode:  Codes.DuplicateUsername,
	},
	AuthenticationFailed: types.ServiceError{
        StatusCode: http.StatusUnauthorized,
        ErrorCode:  "AuthenticationFailed",
    },
}