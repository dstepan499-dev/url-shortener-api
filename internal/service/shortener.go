package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"

	"github.com/dstepan499-dev/url-shortener-api/internal/repository"
)

const (
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	aliasLength = 6
)

var (
	ErrInvalidURL  = errors.New("invalid url format, must start with http:// or https://")
	ErrEmptyURL    = errors.New("url cannot be empty")
	ErrAliasExists = repository.ErrAliasExists
	ErrNotFound    = repository.ErrNotFound
)

type ShortenerService struct {
	repo *repository.Repository
}

func NewShortenerService(repo *repository.Repository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) ShortenURl(ctx context.Context, originalURL, customAlias string) (string, error) {
	originalURL = strings.TrimSpace(originalURL)
	if originalURL == "" {
		return "", ErrEmptyURL
	}

	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return "", ErrInvalidURL
	}

	alias := strings.TrimSpace(customAlias)
	if alias == "" {
		alias, err = generateRandomAlias(aliasLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate alias: %w", err)
		}
	}

	if err := s.repo.SaveURL(ctx, alias, originalURL); err != nil {
		return "", err
	}

	return alias, nil
}

func (s *ShortenerService) ResolveURL(ctx context.Context, alias string) (string, error) {
	return s.repo.GetAndIncrement(ctx, alias)
}

func (s *ShortenerService) GetAnalytics(ctx context.Context, alias string) (*repository.URlRecord, error) {
	return s.repo.GetAnalytics(ctx, alias)
}

func generateRandomAlias(length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
