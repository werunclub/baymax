package errors

import (
	"net/http"
	"testing"

	"github.com/pborman/uuid"
)

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
		ne := New(e.Id, e.Code, e.Detail, e.Status)

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
}
