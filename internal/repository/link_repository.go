package repository

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

type LinkRepository interface {
	CreateLink(link *models.Link) error
	GetLinkByShortCode(shortCode string) (*models.Link, error)
	GetAllLinks() ([]models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}

type GormLinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{
		db: db,
	}
}

func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	result := r.db.Create(link)
	return result.Error
}

func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	result := r.db.Where("short_code = ?", shortCode).First(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	result := r.db.Find(&links)
	return links, result.Error
}

func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	result := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}
