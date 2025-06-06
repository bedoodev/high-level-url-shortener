package repository

import (
	"context"
	"time"

	"github.com/bedoodev/high-level-url-shortener/internal/config"
	"github.com/bedoodev/high-level-url-shortener/internal/model"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	FindByShortCode(ctx context.Context, code string) (*model.URL, error)
	IncrementClickCount(ctx context.Context, code string) error
	GetDailyClickCounts(ctx context.Context, urlID uint) (map[string]int, error)
	GetTopClickedURLs(ctx context.Context, limit int) ([]map[string]interface{}, error)
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

func (r *urlRepo) GetDailyClickCounts(ctx context.Context, urlID uint) (map[string]int, error) {
	rows, err := config.DB.
		Table("click_events").
		Select("DATE(timestamp) AS day, COUNT(*) as count").
		Where("url_id = ?", urlID).
		Group("day").
		Order("day").
		Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[string]int)

	for rows.Next() {
		var day time.Time
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result[day.Format("2006-01-02")] = count
	}
	return result, nil
}

func (r *urlRepo) GetTopClickedURLs(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := config.DB.
		Table("click_events ce").
		Select("u.short_code, COUNT(ce.*) as count").
		Joins("JOIN urls u ON u.id = ce.url_id").
		Group("u.short_code").
		Order("count DESC").
		Limit(limit).
		Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var shortCode string
		var count int
		if err := rows.Scan(&shortCode, &count); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"short_code":  shortCode,
			"click_count": count,
		})
	}
	return results, nil
}
