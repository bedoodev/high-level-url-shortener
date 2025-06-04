package repository

import (
	"context"

	"github.com/bedoodev/high-level-url-shortener/internal/config"
	"github.com/bedoodev/high-level-url-shortener/internal/model"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	FindByShortCode(ctx context.Context, code string) (*model.URL, error)
	IncrementClickCount(ctx context.Context, code string) error
}

type urlRepo struct{}

func NewURLRepository() URLRepository {
	return &urlRepo{}
}

func (r *urlRepo) Create(ctx context.Context, url *model.URL) error {
	return config.DB.WithContext(ctx).Create(url).Error
}

func (r *urlRepo) FindByShortCode(ctx context.Context, code string) (*model.URL, error) {
	var url model.URL
	err := config.DB.WithContext(ctx).Where("short_code = ?", code).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepo) IncrementClickCount(ctx context.Context, code string) error {
	return config.DB.WithContext(ctx).
		Model(&model.URL{}).
		Where("short_code = ?", code).
		UpdateColumn("click_count", gorm.Expr("click_count + 1")).Error
}
