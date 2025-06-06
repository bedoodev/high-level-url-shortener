package service

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/bedoodev/high-level-url-shortener/internal/config"
	"github.com/bedoodev/high-level-url-shortener/internal/kafka"
	"github.com/bedoodev/high-level-url-shortener/internal/model"
	"github.com/bedoodev/high-level-url-shortener/internal/repository"
	"go.uber.org/zap"
)

type URLService interface {
	ShortenURL(ctx context.Context, originalURL string) (*model.URL, error)
	ResolveURL(ctx context.Context, shortCode string) (*model.URL, error)
	GetAnalytics(ctx context.Context, shortCode string) (map[string]int, error)
	GetTopURLs(ctx context.Context, limit int) ([]map[string]interface{}, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) ShortenURL(ctx context.Context, originalURL string) (*model.URL, error) {
	originalURL = strings.TrimSpace(originalURL)

	if originalURL == "" {
		return nil, errors.New("original URL cannot be empty")
	}

	var shortCode string
	var err error

	for i := 0; i < 5; i++ { // max 5 deneme
		shortCode = generateShortCode(7)
		_, err = s.repo.FindByShortCode(ctx, shortCode)
		if err != nil {
			break // kayıt yoksa => unique
		}
	}

	if err == nil {
		return nil, errors.New("could not generate unique short code")
	}

	url := &model.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}

	if err := s.repo.Create(ctx, url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) ResolveURL(ctx context.Context, shortCode string) (*model.URL, error) {
	// 1. Check Local cache (RAM)
	if original, ok := config.HotKeyCache.Get(shortCode); ok {

		zap.L().Info("Response from RAM")
		_ = kafka.PublishClick(ctx, shortCode)
		return &model.URL{
			OriginalURL: original,
			ShortCode:   shortCode,
		}, nil
	}

	// 2. Check redis
	original, err := config.RedisClient.Get(ctx, shortCode).Result()

	if err == nil {
		zap.L().Info("Response from redis")
		_ = kafka.PublishClick(ctx, shortCode)

		// Write to local cache
		config.HotKeyCache.Set(shortCode, original)

		return &model.URL{
			OriginalURL: original,
			ShortCode:   shortCode,
		}, nil

	}

	// 3. Check Postgres
	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	zap.L().Info("Response from db")
	// Write to redis
	_ = config.RedisClient.Set(ctx, shortCode, url.OriginalURL, 24*time.Hour).Err() // 0 => no expiration time, write t

	// Write to local cache
	config.HotKeyCache.Set(shortCode, url.OriginalURL)

	// Update click count
	_ = kafka.PublishClick(ctx, shortCode)

	return url, nil
}

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateShortCode(length int) string {
	code := make([]rune, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func (s *urlService) GetAnalytics(ctx context.Context, shortCode string) (map[string]int, error) {
	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	return s.repo.GetDailyClickCounts(ctx, url.ID)
}

func (s *urlService) GetTopURLs(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	return s.repo.GetTopClickedURLs(ctx, limit)
}
