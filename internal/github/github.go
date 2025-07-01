package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/avast/retry-go/v4"
	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/rs/zerolog"
)

const (
	retryDelaySeconds    int = 5
	maxRetryDelayMinutes int = 2
	attempts                 = 5
)

type GitHubClient interface {
	CloneRepo(ctx context.Context, cloneURL string) (string, error)
}

type Client struct {
	token  string
	logger zerolog.Logger
}

func NewClient(token string, logger zerolog.Logger) *Client {
	return &Client{
		token:  token,
		logger: logger,
	}
}

func (c *Client) CloneRepo(ctx context.Context, cloneURL string) (string, error) {
	val, requestErr := retry.DoWithData(
		func() (string, error) {
			dir, err := os.MkdirTemp("", "grfscan-*")
			if err != nil {
				return "", fmt.Errorf("error creating temp directory: %w", err)
			}

			var auth transport.AuthMethod
			if c.token != "" {
				auth = &http.BasicAuth{
					Password: c.token,
				}
			}

			_, err = git.PlainCloneContext(ctx, dir, &git.CloneOptions{
				URL:      cloneURL,
				Depth:    1,
				Auth:     auth,
				Progress: os.Stdout,
			})
			if err != nil {
				os.RemoveAll(dir)
				return "", fmt.Errorf("error cloning repository: %w", err)
			}

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return dir, nil
			}
			return absDir, nil
		},
		retry.Attempts(attempts),
		retry.OnRetry(func(n uint, err error) {
			c.logger.Error().
				Err(err).
				Uint("retry_attempt", n+1).
				Int("max_retries", attempts).
				Int("initial_delay_seconds", retryDelaySeconds).
				Int("max_delay_minutes", maxRetryDelayMinutes).
				Msg("retrying request to clone repository")
		}),
		retry.Delay(time.Duration(retryDelaySeconds)*time.Second),
		retry.MaxDelay(time.Duration(maxRetryDelayMinutes)*time.Minute),
		retry.DelayType(retry.BackOffDelay),
	)

	if requestErr != nil {
		return "", fmt.Errorf("error cloning repository: %w", requestErr)
	}

	if val == "" {
		return "", fmt.Errorf("cloning repository returned empty directory")
	}

	if _, err := os.Stat(val); os.IsNotExist(err) {
		return "", fmt.Errorf("cloned directory does not exist: %s", val)
	}

	return val, nil
}
