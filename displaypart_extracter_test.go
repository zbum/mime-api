package main

import (
	"net/mail"
	"os"
	"strings"
	"testing"
)

func TestExtractDisplayPart_Multipart(t *testing.T) {
	f, err := os.Open("test/resources/simple.eml")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()

	msg, err := mail.ReadMessage(f)
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	content, contentType, err := ExtractDisplayPart(msg)
	if err != nil {
		t.Fatalf("ExtractDisplayPart failed: %v", err)
	}

	if contentType != "text/html" {
		t.Errorf("expected content type 'text/html', got '%s'", contentType)
	}

	if !strings.Contains(content, "<p><b>hi!</b></p>") {
		t.Errorf("expected content to contain '<p><b>hi!</b></p>', got '%s'", content)
	}
}

func TestExtractDisplayPart_Singlepart(t *testing.T) {
	f, err := os.Open("test/resources/singlepart.eml")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()

	msg, err := mail.ReadMessage(f)
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	content, contentType, err := ExtractDisplayPart(msg)
	if err != nil {
		t.Fatalf("ExtractDisplayPart failed: %v", err)
	}

	if contentType != "text/html" {
		t.Errorf("expected content type 'text/html', got '%s'", contentType)
	}

	if !strings.Contains(content, "<table") {
		t.Errorf("expected content to contain '<table', got truncated content")
	}
}

func TestDecodeContent_Base64(t *testing.T) {
	input := strings.NewReader("SGVsbG8gV29ybGQ=")
	content, err := decodeContent(input, "base64")
	if err != nil {
		t.Fatalf("decodeContent failed: %v", err)
	}

	if content != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", content)
	}
}

func TestDecodeContent_QuotedPrintable(t *testing.T) {
	input := strings.NewReader("Hello=20World")
	content, err := decodeContent(input, "quoted-printable")
	if err != nil {
		t.Fatalf("decodeContent failed: %v", err)
	}

	if content != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", content)
	}
}

func TestDecodeContent_Plain(t *testing.T) {
	input := strings.NewReader("Hello World")
	content, err := decodeContent(input, "")
	if err != nil {
		t.Fatalf("decodeContent failed: %v", err)
	}

	if content != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", content)
	}
}
