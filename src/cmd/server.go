package cmd

import (
	"log"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Webserver CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.InitCfg()
		c, err := cfg.Initiate("api")
		if err != nil {
			log.Fatalf("init setup failed: %v", err)
		}

		// Schema migration dikelola via cmd/migrate.go — server tidak menjalankan AutoMigrate.
		// RabbitMQ reconnect + channel pool dikelola oleh RabbitManager (src/infrastructure/messaging/rabbitmq.go).

		router := http.MapUrls(http.App3rdParty{
			Config: cfg,
			Logs:   c.Logs,
			DB:     c.DB,
			R:      c.R,
			RCP:    c.RCP,
			RM:     c.RM,
			GS:     c.GS,
		})

		log.Fatal(router.Listen(":" + c.Config.AppApiPort))
	},
}
