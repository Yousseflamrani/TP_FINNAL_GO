package repository

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

type ClickRepository interface {
	CreateClick(click *models.Click) error
	CountClicksByLinkID(linkID uint) (int, error)
}

type GormClickRepository struct {
	db *gorm.DB
}

// C'est la méthode recommandée pour obtenir un dépôt, garantissant que la connexion à la base de données est injectée.
func NewClickRepository(db *gorm.DB) *GormClickRepository {
	return &GormClickRepository{db: db}
}

func (r *GormClickRepository) CreateClick(click *models.Click) error {
	if err := r.db.Create(click).Error; err != nil {
		return fmt.Errorf("erreur lors de la création du clic: %w", err)
	}
	return nil
}

func (r *GormClickRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	if err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("erreur lors du comptage des clics pour le lien %d: %w", linkID, err)
	}

	return int(count), nil
}
