package version

import (
	"bytes"
	"fmt"
	"time"
)

// Info contains the version info
type Info struct {
	Revision          string    `json:"revision"`
	Author            string    `json:"author"`
	Version           string    `json:"version"`
	VersionPrerelease string    `json:"version_prerelease"`
	VersionMetadata   string    `json:"version_metadata"`
	BuildDate         string    `json:"build_date"`
	BuildUser         string    `json:"build_user"`
	Branch            string    `json:"branch"`
	BuildHost         string    `json:"build_host"`
	ServerStartTime   time.Time `json:"start_time"`
}

// New creates a new VersionInfo
// These values should be passed in from the main as ldflags
func New(
	v,
	gitCommit,
	gitBranch,
	gitAuthor,
	gitDescribe,
	buildDate,
	versionPrerelease,
	versionMetaData,
	buildUser,
	buildHost string,
	startTime time.Time) Info {

	ver := v
	rel := versionPrerelease

	if gitDescribe != "" {
		ver = gitDescribe
	}
	if gitDescribe == "" && rel == "" {
		rel = "dev"
	}

	return Info{
		Revision:          gitCommit,
		Author:            gitAuthor,
		Version:           ver,
		VersionPrerelease: rel,
		VersionMetadata:   versionMetaData,
		BuildDate:         buildDate,
		BuildUser:         buildUser,
		Branch:            gitBranch,
		BuildHost:         buildHost,
		ServerStartTime:   startTime,
	}
}

// UpTime returns the time.Since(startTime)
func (c *Info) UpTime() time.Duration {
	return time.Since(c.ServerStartTime)
}

// Number returns the string representation of the build number
func (c *Info) Number() string {
	if c.Version == "" && c.VersionPrerelease == "" {
		return "(version unknown)"
	}

	version := fmt.Sprintf("%s", c.Version)

	if c.VersionPrerelease != "" {
		version = fmt.Sprintf("%s-%s", version, c.VersionPrerelease)
	}
	if c.VersionMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, c.VersionMetadata)
	}
	return version
}

// FullVersionNumber returns the Full Version Number
func (c *Info) FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer

	if c.Version == "" && c.VersionPrerelease == "" {
		return "(version unknown)"
	}

	fmt.Fprintf(&versionString, "EveBot %s", c.Version)

	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)
	}

	if c.VersionMetadata != "" {
		fmt.Fprintf(&versionString, "+%s", c.VersionMetadata)
	}

	if rev && c.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	return versionString.String()
}
