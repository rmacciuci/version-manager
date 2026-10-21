// Package semver finds the latest vX.Y.Z tag (optionally under a prefix such
// as "api/") and computes the next one.
package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Bump selects which part of the version to increase. When several are set,
// they apply in order major, minor, patch (major+patch gives X+1.0.1).
type Bump struct {
	Major bool
	Minor bool
	Patch bool
}

// ErrNoBump is returned when no part of the version was selected.
var ErrNoBump = errors.New("no version increase selected: set at least one of major, minor or patch")

var tagPattern = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)

// Version is a parsed vX.Y.Z version.
type Version struct {
	Major, Minor, Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) less(o Version) bool {
	if v.Major != o.Major {
		return v.Major < o.Major
	}
	if v.Minor != o.Minor {
		return v.Minor < o.Minor
	}

	return v.Patch < o.Patch
}

// Parse parses "vX.Y.Z" (no prefix).
func Parse(s string) (Version, bool) {
	m := tagPattern.FindStringSubmatch(s)
	if m == nil {
		return Version{}, false
	}

	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])

	return Version{Major: major, Minor: minor, Patch: patch}, true
}

// Latest returns the highest version among tags that are exactly
// prefix+"vX.Y.Z". Tags with another prefix, or with none when a prefix is
// given, are ignored. found is false when no tag matches (start from v0.0.0).
func Latest(tags []string, prefix string) (latest Version, found bool) {
	for _, tag := range tags {
		rest, ok := strings.CutPrefix(strings.TrimSpace(tag), prefix)
		if !ok {
			continue
		}

		v, ok := Parse(rest)
		if !ok {
			continue
		}

		if !found || latest.less(v) {
			latest, found = v, true
		}
	}

	return latest, found
}

// Next applies the bump to v.
func Next(v Version, b Bump) (Version, error) {
	if !b.Major && !b.Minor && !b.Patch {
		return Version{}, ErrNoBump
	}

	if b.Major {
		v = Version{Major: v.Major + 1}
	}
	if b.Minor {
		v = Version{Major: v.Major, Minor: v.Minor + 1}
	}
	if b.Patch {
		v.Patch++
	}

	return v, nil
}
