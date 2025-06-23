package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/axellelanca/urlshortener/cmd" // Pour accéder à cmd.Cfg
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var longURLFlag string // Flag pour l'URL longue à raccourcir

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		if longURLFlag == "" {
			fmt.Println("ERREUR: Le flag --url est requis.")
			os.Exit(1)
		}

		_, err := url.ParseRequestURI(longURLFlag)
		if err != nil {
			fmt.Println("ERREUR: L'URL fournie n'est pas valide.")
			os.Exit(1)
		}

		cfg := cmd.Cfg // <- récupère la config chargée globalement

		// Connexion BDD SQLite
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erreur ouverture BDD: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Erreur obtention instance SQL: %v", err)
		}
		defer sqlDB.Close()

		// Créer le service
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// Créer le lien court
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			log.Fatalf("Erreur création lien court: %v", err)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Println("✅ URL courte créée avec succès:")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("Lien complet: %s\n", fullShortURL)
	},
}

func init() {
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")
	CreateCmd.MarkFlagRequired("url")
	cmd.RootCmd.AddCommand(CreateCmd)
}
