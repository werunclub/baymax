package errors

import (
	"encoding/json"
	"net/http"

	"github.com/pborman/uuid"
)

// Errors provide a way to return detailed information
// for an RPC request error. The error is normally
// JSON encoded.
type Error struct {
	Id     string `json:"id"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
	Status int32  `json:"status"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func New(id string, code string, detail string, status int32) error {
	return &Error{
		Id:     id,
		Code:   code,
		Detail: detail,
		Status: status,
	}
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

func BadRequest(detail string) error {
	return &Error{
		Id:     uuid.New(),
		Code:   "bad_request",
		Detail: detail,
		Status: http.StatusBadRequest,
	}
}

func Unauthorized(detail string) error {
	return &Error{
		Id:     uuid.New(),
		Code:   "unauthorized",
		Detail: detail,
		Status: http.StatusUnauthorized,
	}
}

func Forbidden(detail string) error {
	return &Error{
		Id:     uuid.New(),
		Code:   "forbidden",
		Detail: detail,
		Status: http.StatusForbidden,
	}
}

func NotFound(detail string) error {
	return &Error{
		Id:     uuid.New(),
		Code:   "not_found",
		Detail: detail,
		Status: http.StatusNotFound,
	}
}

func InternalServerError(detail string) error {
	return &Error{
		Id:     uuid.New(),
		Code:   "internal_server_error",
		Detail: detail,
		Status: http.StatusInternalServerError,
	}
}
