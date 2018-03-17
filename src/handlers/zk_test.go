package handlers

import (
	"net/http"
	"testing"
)

func TestList(t *testing.T) {
	req1, err := http.NewRequest("GET", "/zk/106", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := newRequestRecorder(req1, "GET", "/zk/:cluster", List)
	if rr1.Code != 404 {
		t.Error("Expected response code to be 404")
	}
	// expected response
	er1 := "{\"error\":{\"status\":404,\"title\":\"Record Not Found\"}}\n"
	if rr1.Body.String() != er1 {
		t.Error("Response body does not match")
	}

	t.Log("When the book exists")
	t.Log(rr1.Body.String())
}
