package download

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
)

func LatestTronVersion() (string, error) {
	const url = "https://api.github.com/repos/loghole/tron/releases/latest"

	resp, err := http.Get(url) // nolint:gosec,bodyclose,noctx //body is closed
	if err != nil {
		return "", fmt.Errorf("http get '%s': %w", url, err)
	}

	defer helpers.Close(resp.Body)

	var dest models.Release

	if err := jsoniter.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return "", fmt.Errorf("decode json body: %w", err)
	}

	return dest.TagName, nil
}
