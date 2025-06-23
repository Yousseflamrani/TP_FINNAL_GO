package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type LinkService struct {
	linkRepo repository.LinkRepository
}

func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

func (s *LinkService) GenerateShortCode(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("erreur lors de la génération du code court: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	var shortCode string
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		code, err := s.GenerateShortCode(6)

		if err != nil {
			return nil, fmt.Errorf("erreur lors de la la génération du code court: %w", err)
		}

		_, err = s.linkRepo.GetLinkByShortCode(code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code
				break
			}
			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}

		log.Printf("Short code '%s' already exists, retrying generation (%d/%d)...", code, i+1, maxRetries)
	}

	if shortCode == "" {
		return nil, errors.New("impossible de générer un code court unique après plusieurs tentatives")
	}

	link := &models.Link{
		LongURL:   longURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
	}

	err := s.linkRepo.CreateLink(link)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du lien: %w", err)
	}

	return link, nil
}

func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du lien: %w", err)
	}
	return link, nil
}

func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors de la récupération du lien: %w", err)
	}

	clickCount, err := s.linkRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors du comptage des clics: %w", err)
	}

	return link, clickCount, nil
}
