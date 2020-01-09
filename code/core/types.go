package core

import (
	"errors"
	"strconv"
	"strings"
)

type AppVersion struct {
	MajorVersion int64
	MinorVersion int64
}

func (a *AppVersion)String() string {
	return strconv.FormatInt(a.MajorVersion, 10) + "." + strconv.FormatInt(a.MinorVersion, 10)
}

func GetAppVersionFromString(s string) (*AppVersion, error) {
	tokens := strings.Split(s, ".")
	numTokens := len(tokens)

	if numTokens < 1 || numTokens > 2 {
		return nil, errors.New("malformed app version string")
	}

	majorVersion, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return nil, errors.New("error parsing major version: " + err.Error())
	}

	minorVersion := int64(0)
	if numTokens > 1 {
		minorVersion, err = strconv.ParseInt(tokens[1], 10, 64)
		if err != nil {
			return nil, errors.New("error parsing minor version: " + err.Error())
		}
	}

	return &AppVersion{MajorVersion: majorVersion, MinorVersion: minorVersion}, nil
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
