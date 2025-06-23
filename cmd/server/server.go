package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/api"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/monitor"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/axellelanca/urlshortener/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RunServerCmd représente la commande 'run-server' de Cobra.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatal("Configuration non chargée")
		}

		// Initialiser la connexion à la base de données SQLite avec GORM.
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erreur lors de la connexion à la base de données: %v", err)
		}

		// Migration automatique des modèles
		err = db.AutoMigrate(&models.Link{}, &models.Click{})
		if err != nil {
			log.Fatalf("Erreur lors de la migration de la base de données: %v", err)
		}

		// Initialiser les repositories.
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)
		log.Println("Repositories initialisés.")

		// Initialiser les services métiers.
		linkService := services.NewLinkService(linkRepo)
		log.Println("Services métiers initialisés.")

		// Initialiser le channel ClickEventsChannel (api/handlers) des événements de clic et lancer les workers (StartClickWorkers).
		api.ClickEventsChannel = make(chan models.ClickEvent, cfg.Analytics.BufferSize)
		workers.StartClickWorkers(3, api.ClickEventsChannel, clickRepo) // 3 workers
		log.Printf("Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
			cfg.Analytics.BufferSize, 3)

		// Initialiser et lancer le moniteur d'URLs.
		monitorInterval := time.Duration(cfg.Monitor.IntervalMinutes) * time.Minute
		urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)
		go urlMonitor.Start()
		log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)

		// Configurer le routeur Gin et les handlers API.
		router := gin.Default()
		api.SetupRoutes(router, linkService)

		log.Println("Routes API configurées.")

		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		go func() {
			log.Printf("Serveur HTTP démarré sur le port %d", cfg.Server.Port)
			log.Printf("URL de base: %s", cfg.Server.BaseURL)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Erreur du serveur HTTP: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")

		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
