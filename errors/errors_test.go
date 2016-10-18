package errors

import (
	"net/http"
	"testing"

	"github.com/pborman/uuid"
)

func getNilError() *Error {
	return nil
}

func getError() *Error {
	return Parse(BadRequest("not_found").Error())
}

func TestErrors(t *testing.T) {
	testData := []*Error{
		&Error{
			Id:     uuid.New(),
			Code:   "status_internal_server_error",
			Detail: "Internal server error",
			Status: http.StatusInternalServerError,
		},
	}

	for _, e := range testData {
		ne := New(e.Code, e.Detail, e.Status)

		pe := Parse(ne.Error())

		if pe == nil {
			t.Fatalf("Expected error got nil %v", pe)
		}

		if pe.Detail != e.Detail {
			t.Fatalf("Expected %s got %s", e.Detail, pe.Detail)
		}

		if pe.Code != e.Code {
			t.Fatalf("Expected %s got %s", e.Code, pe.Code)
		}

		if pe.Status != e.Status {
			t.Fatalf("Expected %s got %s", e.Status, pe.Status)
		}
	}

	if err := getNilError(); err != nil {
		t.Fatalf("Expected nil got %s", err)
	}

	if err := getError(); err == nil {
		t.Fatal("Expected error got nil")
	}
}
