package workers

import (
	"log"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

// StartClickWorkers lance un pool de goroutines "workers" pour traiter les événements de clic.
// Chaque worker lira depuis le même 'clickEventsChan' et utilisera le 'clickRepo' pour la persistance.
func StartClickWorkers(workerCount int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	log.Printf("Starting %d click worker(s)...", workerCount)
	for i := 0; i < workerCount; i++ {
		go clickWorker(clickEventsChan, clickRepo)
	}
}

// clickWorker est la fonction exécutée par chaque goroutine worker.
func clickWorker(clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	for event := range clickEventsChan {
		click := &models.Click{
			LinkID:    event.LinkID,
			UserAgent: event.UserAgent,
		}

		err := clickRepo.CreateClick(click)

		if err != nil {
			log.Printf("ERROR: Failed to save click for LinkID %d (UserAgent: %s): %v",
				event.LinkID, event.UserAgent, err)
			log.Printf("Click recorded successfully for LinkID %d", event.LinkID)
		}
	}
}
