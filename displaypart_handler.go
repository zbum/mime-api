package main

import (
	"net/http"
	"net/mail"
)

func DisplayPart(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}
	defer func() { _ = f.Close() }()

	message, err := mail.ReadMessage(f)
	if err != nil {
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	content, contentType, err := ExtractDisplayPart(message)
	if err != nil {
		http.Error(w, "parse failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write([]byte(content))
}