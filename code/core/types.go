package core

import (
	"strconv"
)

type AppVersion struct {
	MajorVersion int64
	MinorVersion int64
}

func (a *AppVersion)String() string {
	return strconv.FormatInt(a.MajorVersion, 10) + "." + strconv.FormatInt(a.MinorVersion, 10)
}

func (a *AppVersion)Equals(b *AppVersion) bool {
	return a.MajorVersion == b.MajorVersion &&
		   a.MinorVersion == b.MinorVersion
}

func (a *AppVersion)Compare(b *AppVersion) int {
	if a.MajorVersion < b.MajorVersion {
		return -1
	}
	if a.MajorVersion > b.MajorVersion {
		return 1
	}

	// Major versions are equal

	if a.MinorVersion < b.MinorVersion {
		return -1
	}
	if a.MinorVersion > b.MinorVersion {
		return 1
	}

	// Major and Minor versions are equal

	return 0
}
