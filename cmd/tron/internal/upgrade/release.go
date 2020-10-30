package upgrade

import (
	"time"

	"github.com/Masterminds/semver"
)

type release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
}

func (r *release) version() *semver.Version {
	v, err := semver.NewVersion(r.TagName)
	if err != nil {
		return nil
	}

	return v
}
