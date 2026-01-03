package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDisplayPart_Success(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open("test/resources/simple.eml")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "simple.eml")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("failed to copy file: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/v1/display-part", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	DisplayPart(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("expected content type 'text/html', got '%s'", w.Header().Get("Content-Type"))
	}

	if !strings.Contains(w.Body.String(), "<p><b>hi!</b></p>") {
		t.Errorf("expected body to contain '<p><b>hi!</b></p>'")
	}
}

func TestDisplayPart_NoFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/v1/display-part", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	DisplayPart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDisplayPart_InvalidMultipart(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/display-part", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	DisplayPart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}