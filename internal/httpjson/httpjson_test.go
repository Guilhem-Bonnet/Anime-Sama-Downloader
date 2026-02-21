package httpjson

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrite_SetsContentTypeAndStatus(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, http.StatusCreated, map[string]string{"hello": "world"})

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	expected := "application/json; charset=utf-8"
	if ct != expected {
		t.Errorf("expected Content-Type %q, got %q", expected, ct)
	}
}

func TestWrite_EncodesJSON(t *testing.T) {
	w := httptest.NewRecorder()

	payload := map[string]int{"count": 42}
	Write(w, http.StatusOK, payload)

	var got map[string]int
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got["count"] != 42 {
		t.Errorf("expected count=42, got %d", got["count"])
	}
}

func TestWrite_NilPayload(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, http.StatusOK, nil)

	body := w.Body.String()
	if body != "null\n" {
		t.Errorf("expected null JSON, got %q", body)
	}
}

func TestWrite_StructPayload(t *testing.T) {
	type resp struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	w := httptest.NewRecorder()
	Write(w, http.StatusOK, resp{Name: "Gopher", Age: 15})

	var got resp
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if got.Name != "Gopher" || got.Age != 15 {
		t.Errorf("unexpected payload: %+v", got)
	}
}

func TestWriteError_FormatsErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()

	WriteError(w, http.StatusBadRequest, "invalid input")

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var got ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if got.Error != "invalid input" {
		t.Errorf("expected error message %q, got %q", "invalid input", got.Error)
	}
}

func TestWriteError_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()

	WriteError(w, http.StatusInternalServerError, "something broke")

	var got ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if got.Error != "something broke" {
		t.Errorf("expected error %q, got %q", "something broke", got.Error)
	}
}

func TestWrite_EmptySlice(t *testing.T) {
	w := httptest.NewRecorder()

	Write(w, http.StatusOK, []string{})

	body := w.Body.String()
	if body != "[]\n" {
		t.Errorf("expected empty array, got %q", body)
	}
}
