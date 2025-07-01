package output

import (
	"encoding/json"
	"fmt"
	"os"

	"gr-scanner/internal/models"
)

// Writer handles output generation
type Writer struct{}

// New creates a new Writer
func New() *Writer {
	return &Writer{}
}

// Write prints the result as JSON to stdout
func (w *Writer) Write(result *models.Output) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding output JSON: %w", err)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		return fmt.Errorf("writing to stdout: %w", err)
	}
	return nil
}
