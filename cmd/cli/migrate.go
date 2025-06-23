package cli

import (
	"fmt"
	"log"

	"github.com/axellelanca/urlshortener/cmd" // pour cmd.RootCmd
	"github.com/axellelanca/urlshortener/config"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Charger la configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Erreur chargement configuration: %v", err)
		}

		// Connexion à SQLite
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erreur connexion base SQLite: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Erreur accès instance SQL: %v", err)
		}
		defer sqlDB.Close()

		// Migration des modèles
		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("Erreur migration GORM: %v", err)
		}

		fmt.Println("✅ Migrations exécutées avec succès.")
	},
}

func init() {
	cmd.RootCmd.AddCommand(MigrateCmd)
}
