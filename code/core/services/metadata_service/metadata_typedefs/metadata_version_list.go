package metadata_typedefs

import (
	"errors"

	"github.com/spacetimi/timi_shared_server/v2/code/core"
)

type MetadataVersionList struct {
	Versions        []string
	CurrentVersions []string

	_versionsAsMap map[string]bool
}

func (mvl *MetadataVersionList) Initialize() {
	mvl._versionsAsMap = make(map[string]bool)
	for _, version := range mvl.Versions {
		mvl._versionsAsMap[version] = true
	}
}

func (mvl *MetadataVersionList) IsVersionValid(version *core.AppVersion) bool {
	_, ok := mvl._versionsAsMap[version.String()]
	return ok
}

func (mvl *MetadataVersionList) IsVersionCurrent(version *core.AppVersion) bool {
	for _, currentVersion := range mvl.CurrentVersions {
		if version.String() == currentVersion {
			return true
		}
	}
	return false
}

func (mvl *MetadataVersionList) CreateNewVersion(version *core.AppVersion, markAsCurrent bool) error {
	_, ok := mvl._versionsAsMap[version.String()]
	if ok {
		return errors.New("duplicate version")
	}

	mvl.Versions = append(mvl.Versions, version.String())
	mvl._versionsAsMap[version.String()] = true

	if markAsCurrent {
		mvl.CurrentVersions = append(mvl.CurrentVersions, version.String())
	}

	return nil
}

func (mvl *MetadataVersionList) GetLatestVersionDefined() (*core.AppVersion, error) {
	var latestVersion *core.AppVersion
	for _, versionString := range mvl.Versions {
		version, err := core.GetAppVersionFromString(versionString)
		if err != nil {
			continue
		}
		if latestVersion == nil ||
			latestVersion.Compare(version) < 0 {
			latestVersion = version
		}
	}

	if latestVersion == nil {
		return nil, errors.New("couldn't find latest version")
	}
	return latestVersion, nil
}
