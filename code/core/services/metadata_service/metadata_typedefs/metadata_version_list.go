package metadata_typedefs

import "github.com/spacetimi/server/timi_shared/code/core"

type MetadataVersionList struct {
	Versions []string
	CurrentVersions []string

	_versionsAsMap map[string]bool
}

func (mvl *MetadataVersionList) Initialize() {
	mvl._versionsAsMap = make(map[string]bool)
	for _, version := range mvl.Versions {
		mvl._versionsAsMap[version] = true
	}
}

func (mvl *MetadataVersionList)IsVersionValid(version *core.AppVersion) bool {
	_, ok := mvl._versionsAsMap[version.String()]
	return ok
}
func (mvl *MetadataVersionList)IsVersionCurrent(version *core.AppVersion) bool {
	for _, currentVersion := range mvl.CurrentVersions {
		if version.String() == currentVersion {
			return true
		}
	}
	return false
}

