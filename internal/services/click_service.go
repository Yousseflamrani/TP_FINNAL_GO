package services

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

type ClickService struct {
	clickRepo repository.ClickRepository
}

func NewClickService(clickRepo repository.ClickRepository) *ClickService {
	return &ClickService{
		clickRepo: clickRepo,
	}
}

func (s *ClickService) RecordClick(click *models.Click) error {
	err := s.clickRepo.CreateClick(click)
	if err != nil {
		return fmt.Errorf("erreur lors de l'enregistrement du clic: %w", err)
	}
	return nil
}

func (s *ClickService) GetClicksCountByLinkID(linkID uint) (int, error) {
	count, err := s.clickRepo.CountClicksByLinkID(linkID)
	if err != nil {
		return 0, fmt.Errorf("erreur lors du comptage des clics pour le lien ID %d: %w", linkID, err)
	}
	return count, nil
}
