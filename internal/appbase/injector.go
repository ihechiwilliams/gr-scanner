package appbase

import (
	"os"

	"gr-scanner/internal/github"
	"gr-scanner/internal/output"
	"gr-scanner/internal/parser"
	"gr-scanner/internal/scanner"
	"gr-scanner/internal/service"

	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func NewInjector(serviceName string, cfg *Config) *do.Injector {
	injector := do.New()

	// ===========================
	//	Service Configs (logging, open-api,...)
	// ===========================
	do.Provide(injector, func(i *do.Injector) (*zerolog.Logger, error) {
		logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			return nil, err
		}

		logger := zerolog.New(os.Stdout).
			Level(logLevel).
			With().
			Str("serviceName", serviceName).
			Logger()

		return &logger, nil
	})

	do.Provide(injector, func(i *do.Injector) (*parser.ConfigParser, error) {
		return parser.New(), nil
	})

	do.Provide(injector, func(i *do.Injector) (*github.Client, error) {
		logger := do.MustInvoke[*zerolog.Logger](i)
		return github.NewClient(
			cfg.GitHubToken,
			*logger,
		), nil
	})

	do.Provide(injector, func(i *do.Injector) (*scanner.Scanner, error) {
		logger := do.MustInvoke[*zerolog.Logger](i)
		return scanner.New(*logger), nil
	})

	do.Provide(injector, func(i *do.Injector) (*output.Writer, error) {
		return output.New(), nil
	})

	do.Provide(injector, func(i *do.Injector) (*service.Service, error) {
		logger := do.MustInvoke[*zerolog.Logger](i)

		return service.NewService(
			do.MustInvoke[*parser.ConfigParser](i),
			do.MustInvoke[*github.Client](i),
			do.MustInvoke[*scanner.Scanner](i),
			do.MustInvoke[*output.Writer](i),
			*logger,
		), nil
	})

	return injector
}
