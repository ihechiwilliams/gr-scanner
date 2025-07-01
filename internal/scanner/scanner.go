package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"gr-scanner/internal/models"

	"github.com/rs/zerolog"
)

// Scanner handles file scanning
type Scanner struct {
	logger zerolog.Logger
}

// New creates a new Scanner
func New(logger zerolog.Logger) *Scanner {
	return &Scanner{
		logger: logger,
	}
}

// Scan traverses the directory and finds files larger than the threshold
func (s *Scanner) Scan(root string, sizeThreshold int64) (*models.Output, error) {
	var files []models.FileInfo
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // continue
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil // continue
		}

		s.logger.Info().
			Str("path", path).
			Str("size_threshold", fmt.Sprint(sizeThreshold)).
			Int64("size", info.Size()).
			Msg("scanning file")

		if info.Size() > sizeThreshold {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return fmt.Errorf("getting relative path for %s: %w", path, err)
			}

			s.logger.Info().
				Str("path", relPath).
				Int64("size", info.Size()).
				Msg("found large file")

			files = append(files, models.FileInfo{
				Name: relPath,
				Size: info.Size(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scanning directory: %w", err)
	}

	return &models.Output{
		Total: len(files),
		Files: files,
	}, nil
}
