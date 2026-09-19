package service

import (
	"strings"
	"testing"
)

func TestGenerateRandomAlias(t *testing.T) {
	lengths := []int{4, 6, 8, 12}

	for _, length := range lengths {
		alias, err := generateRandomAlias(length)
		if err != nil {
			t.Fatalf("unexpected error for length %d: %v", length, err)
		}

		if len(alias) != length {
			t.Errorf("expected alias length %d, got %d", length, len(alias))
		}

		for _, char := range alias {
			if !strings.ContainsRune(charset, char) {
				t.Errorf("alias contains invalid character: %c", char)
			}
		}
	}
}

func TestShortenURL_Validation(t *testing.T) {
	service := &ShortenerService{}

	tests := []struct {
		name        string
		inputURL    string
		expectedErr error
	}{
		{
			name:        "empty url",
			inputURL:    "",
			expectedErr: ErrEmptyURL,
		},
		{
			name:        "url without scheme",
			inputURL:    "google.com",
			expectedErr: ErrInvalidURL,
		},
		{
			name:        "unsupported scheme ftp",
			inputURL:    "ftp://example.com/file",
			expectedErr: ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ShortenURL(nil, tt.inputURL, "")
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
