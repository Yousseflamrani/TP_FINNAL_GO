package cmd

import (
	"log"

	"github.com/axellelanca/urlshortener/config"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "url-shortener",
		Short: "Service de raccourcissement d'URL avec suivi de clics.",
	}
	Cfg *config.Config // ← Variable globale à utiliser dans les autres fichiers
)

func Execute() {
	var err error
	Cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Impossible de charger la configuration: %v", err)
	}
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("FATAL: Erreur lors de l'exécution de la commande: %v", err)
	}
}
