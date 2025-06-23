package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/config"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var shortCodeFlag string

var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		if shortCodeFlag == "" {
			fmt.Println("❌ Le flag --code est requis.")
			os.Exit(1)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Erreur chargement config: %v", err)
		}

		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erreur connexion base SQLite: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Erreur accès SQL: %v", err)
		}
		defer sqlDB.Close()

		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err == gorm.ErrRecordNotFound {
			fmt.Println("❌ Aucun lien trouvé pour ce code.")
			os.Exit(1)
		} else if err != nil {
			log.Fatalf("Erreur récupération stats: %v", err)
		}

		fmt.Printf("📊 Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("🔗 URL longue: %s\n", link.LongURL)
		fmt.Printf("👁️ Total de clics: %d\n", totalClicks)
	},
}

func init() {
	StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Code court de l'URL (obligatoire)")
	StatsCmd.MarkFlagRequired("code")

	cmd.RootCmd.AddCommand(StatsCmd)
}
