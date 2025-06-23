package repository

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// LinkRepository définit les opérations disponibles sur les liens
type LinkRepository interface {
	CreateLink(link *models.Link) error
	GetLinkByShortCode(shortCode string) (*models.Link, error)
	GetAllLinks() ([]models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}

type GormLinkRepository struct {
	db *gorm.DB
}

// Constructeur du repository
func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{db: db}
}

// Crée un lien dans la base de données
func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	return r.db.Create(link).Error
}

// Récupère un lien par son code court
func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	err := r.db.Where("short_code = ?", shortCode).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// Récupère tous les liens
func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	err := r.db.Find(&links).Error
	if err != nil {
		return nil, err
	}
	return links, nil
}

// Compte le nombre de clics associés à un lien
func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
