package service

import (
	"context"

	"gr-scanner/internal/github"
	"gr-scanner/internal/output"
	"gr-scanner/internal/parser"
	"gr-scanner/internal/scanner"

	"github.com/rs/zerolog"
)

type Service struct {
	config  *parser.ConfigParser
	github  github.GitHubClient
	scanner *scanner.Scanner
	output  *output.Writer
	logger  zerolog.Logger
}

// NewService creates a new Service instance
func NewService(config *parser.ConfigParser, ghClient github.GitHubClient, scanner *scanner.Scanner, outPut *output.Writer, logger zerolog.Logger) *Service {
	return &Service{
		config:  config,
		github:  ghClient,
		scanner: scanner,
		output:  outPut,
		logger:  logger,
	}
}

func (s *Service) Scan(jsonStr string) error {
	ctx := context.Background()

	cfg, err := s.config.Parse(jsonStr)
	if err != nil {
		return err
	}

	// Log the parsed configuration
	s.logger.Info().
		Str("clone_url", cfg.CloneURL).
		Float64("size_mb", cfg.Size).
		Msg("Starting scan with configuration")

	cloneDir, err := s.github.CloneRepo(ctx, cfg.CloneURL)
	if err != nil {
		return err
	}

	s.logger.Info().
		Str("clone_url", cfg.CloneURL).
		Str("clone_dir", cloneDir).
		Msg("downloaded repository")

	sizeThreshold := int64(cfg.Size * 1024 * 1024)
	result, err := s.scanner.Scan(cloneDir, sizeThreshold)
	if err != nil {
		return err
	}

	s.logger.Info().
		Int("total_files", result.Total).
		Msg("scan finished")

	return s.output.Write(result)
}
