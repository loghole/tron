package models

import (
	"time"

	"github.com/Masterminds/semver"
)

type Release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
}

func (r *Release) Version() *semver.Version {
	v, err := semver.NewVersion(r.TagName)
	if err != nil {
		return nil
	}

	return v
}
