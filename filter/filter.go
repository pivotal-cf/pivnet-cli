package filter

import (
	"fmt"
	"path/filepath"
	"regexp"

	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logger"
)

type Filter struct {
	l logger.Logger
}

func NewFilter(l logger.Logger) *Filter {
	return &Filter{
		l: l,
	}
}

// ReleasesByVersion returns all releases that match the provided version regex
func (f Filter) ReleasesByVersion(releases []pivnet.Release, version string) ([]pivnet.Release, error) {
	filteredReleases := make([]pivnet.Release, 0)

	for _, release := range releases {
		match, err := regexp.MatchString(version, release.Version)
		if err != nil {
			return nil, err
		}

		if match {
			filteredReleases = append(filteredReleases, release)
		}
	}

	return filteredReleases, nil
}

func (f Filter) ProductFileNamesByGlobs(
	productFiles []pivnet.ProductFile,
	globs []string,
) ([]pivnet.ProductFile, error) {
	f.l.Debug("filter.ProductFilesNamesByGlobs", logger.Data{"globs": globs})

	filtered := []pivnet.ProductFile{}
	for _, pattern := range globs {
		prevFilteredCount := len(filtered)

		for _, p := range productFiles {
			matched, err := filepath.Match(pattern, p.Name)
			if err != nil {
				return nil, err
			}

			if matched {
				filtered = append(filtered, p)
			}
		}

		if len(filtered) == prevFilteredCount {
			return nil, fmt.Errorf("no product files match glob: %s", pattern)
		}
	}

	return filtered, nil
}
