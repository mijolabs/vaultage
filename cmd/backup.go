package cmd

import (
	"log"

	"github.com/spf13/viper"
)

func performBackup(cfg *viper.Viper) error {
	log.Printf("Performing backup...")
	return nil
}
