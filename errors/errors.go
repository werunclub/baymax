package errors

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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
	if e != nil {
		b, _ := json.Marshal(e)
		return string(b)
	}
	return ""
}

func (e *Error) IsNotFound() bool {
	return e.Code == "not_found"
}

func (e *Error) IsBadRequest() bool {
	return e.Code == "bad_request"
}

func (e *Error) IsUnauthorized() bool {
	return e.Code == "unauthorized"
}

func (e *Error) IsForbidden() bool {
	return e.Code == "forbidden"
}

func (e *Error) IsInternalServerError() bool {
	return e.Code == "internal_server_error"
}

func NewId() string {
	return uuid.New().String()
}

func New(code string, detail string, status int32) error {
	return &Error{
		Id:     NewId(),
		Code:   code,
		Detail: detail,
		Status: status,
	}
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Id = NewId()
		e.Code = "internal_server_error"
		e.Detail = err
		e.Status = 500
	}
	return e
}

func BadRequest(detail string) error {
	return &Error{
		Id:     NewId(),
		Code:   "bad_request",
		Detail: detail,
		Status: http.StatusBadRequest,
	}
}

func Unauthorized(detail string) error {
	return &Error{
		Id:     NewId(),
		Code:   "unauthorized",
		Detail: detail,
		Status: http.StatusUnauthorized,
	}
}

func Forbidden(detail string) error {
	return &Error{
		Id:     NewId(),
		Code:   "forbidden",
		Detail: detail,
		Status: http.StatusForbidden,
	}
}

func NotFound(detail string) error {
	return &Error{
		Id:     NewId(),
		Code:   "not_found",
		Detail: detail,
		Status: http.StatusNotFound,
	}
}

func InternalServerError(detail string) error {
	return &Error{
		Id:     NewId(),
		Code:   "internal_server_error",
		Detail: detail,
		Status: http.StatusInternalServerError,
	}
}
