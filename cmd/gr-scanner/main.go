package main

import (
	"os"

	"gr-scanner/internal/appbase"
	"gr-scanner/internal/service"

	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

const (
	serviceName = "gr-scanner"
)

func main() {
	app := appbase.New(
		appbase.Init(serviceName),
		appbase.WithDependencyInjector(),
	)

	defer app.Shutdown()

	rootCmd := &cobra.Command{
		Use:   "gr-scanner",
		Short: "A CLI tool to scan GitHub repositories for large files",
	}

	scanCmd := &cobra.Command{
		Use:   "scan [json-config]",
		Short: "Scan a repository for files larger than a specified size",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			svc := do.MustInvoke[*service.Service](app.Injector)

			log.Info().Str("command", "scan").Str("config", args[0]).Msg("Starting scan command")

			if err := svc.Scan(args[0]); err != nil {
				log.Error().Err(err).Msg("Scan command failed")
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(scanCmd)
	if err := rootCmd.Execute(); err != nil {
		// log the error using the logger
		log.Error().Err(err).Msg("Command execution failed")
		os.Exit(1)
	}
}
