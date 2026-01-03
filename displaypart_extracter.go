package main

import (
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

func ExtractDisplayPart(message *mail.Message) (string, string, error) {
	mediaType, params, err := mime.ParseMediaType(message.Header.Get("Content-Type"))
	if err != nil {
		return "", "", err
	}

	boundary := params["boundary"]
	if boundary == "" {
		content, err := decodeContent(message.Body, message.Header.Get("Content-Transfer-Encoding"))
		if err != nil {
			return "", "", err
		}
		if strings.HasPrefix(mediaType, "text/html") {
			return content, "text/html", nil
		}
		return content, "text/plain", nil
	}

	reader := multipart.NewReader(message.Body, boundary)

	var textPlain string
	var textHtml string

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", err
		}

		contentType := part.Header.Get("Content-Type")
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			_ = part.Close()
			continue
		}

		transferEncoding := part.Header.Get("Content-Transfer-Encoding")
		content, err := decodeContent(part, transferEncoding)
		_ = part.Close()
		if err != nil {
			continue
		}

		if strings.HasPrefix(mediaType, "text/html") {
			textHtml = content
		} else if strings.HasPrefix(mediaType, "text/plain") {
			textPlain = content
		}
	}

	if textHtml != "" {
		return textHtml, "text/html", nil
	}
	if textPlain != "" {
		return textPlain, "text/plain", nil
	}

	return "", "", nil
}

func decodeContent(r io.Reader, encoding string) (string, error) {
	encoding = strings.ToLower(strings.TrimSpace(encoding))

	var reader io.Reader
	switch encoding {
	case "base64":
		reader = base64.NewDecoder(base64.StdEncoding, r)
	case "quoted-printable":
		reader = quotedprintable.NewReader(r)
	default:
		reader = r
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}