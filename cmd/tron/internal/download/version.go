package download

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
)

func LatestTronVersion() (string, error) {
	const (
		url     = "https://api.github.com/repos/loghole/tron/releases/latest"
		timeout = time.Second * 2
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("http new request '%s': %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req) // nolint:bodyclose // body is closed
	if err != nil {
		return "", fmt.Errorf("http do request '%s': %w", url, err)
	}

	defer helpers.Close(resp.Body)

	var dest models.Release

	if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return "", fmt.Errorf("decode json body: %w", err)
	}

	return dest.TagName, nil
}
