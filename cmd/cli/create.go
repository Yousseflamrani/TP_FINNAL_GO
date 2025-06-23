package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/config"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var longURLFlag string

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		if longURLFlag == "" {
			fmt.Println("Erreur: Le flag --url est requis")
			os.Exit(1)
		}

		_, err := url.ParseRequestURI(longURLFlag)
		if err != nil {
			fmt.Printf("Erreur: Format d'URL invalide: %v\n", err)
			os.Exit(1)
		}
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("FATAL: Impossible de charger la configuration: %v", err)
		}

		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		defer sqlDB.Close()

		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			fmt.Printf("Erreur lors de la création du lien: %v\n", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

func init() {
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "L'URL longue à raccourcir")
	CreateCmd.MarkFlagRequired("url")
	cmd2.RootCmd.AddCommand(CreateCmd)
}
