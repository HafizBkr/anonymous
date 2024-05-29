package utils

import (
	"anonymous/types"
	"net/http"
)

var (
	ErrDecodingBody = types.ServiceError{StatusCode: http.StatusBadRequest, ErrorCode: "DecodingFailed"}
)
