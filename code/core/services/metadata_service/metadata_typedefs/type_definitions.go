package metadata_typedefs

type MetadataSpace int
const (
	METADATA_SPACE_SHARED MetadataSpace = iota
	METADATA_SPACE_APP
)
func (ms MetadataSpace)String() string {
	switch ms {
	case METADATA_SPACE_SHARED:
		return "shared"
	case METADATA_SPACE_APP:
		return "app"
	}
	return "unknown"
}

type IMetadataItem interface {
	GetKey() string
	GetMetadataSpace() MetadataSpace
}

type IMetadataFetcher interface {
	GetMetadataJsonByKey(key string, version string) (string, error)
	GetMetadataVersionList() (*MetadataVersionList, error)
	GetMetadataManifestForVersion(version string) (*MetadataManifest, error)

	/**
     * Only meant to be called from the admin tool / scripts
     */
	SetMetadataVersionList(mvl *MetadataVersionList) error
}

// Error strings
const ERROR_FAILED_TO_READ_METADATA_FILE  = "failed to read metadata file"
const ERROR_FAILED_TO_READ_METADATA_VERSIONS_LIST  = "failed to read metadata versions list"
const ERROR_FAILED_TO_DESERIALIZE_METADATA_VERSIONS_LIST  = "failed to deserialize metadata versions list"
const ERROR_FAILED_TO_READ_METADATA_MANIFEST  = "failed to read metadata manifest"
const ERROR_FAILED_TO_DESERIALIZE_METADATA_MANIFEST  = "failed to deserialize metadata manifest"

type MetadataCache struct {
	Cache map[string]string		// key => metadata json for key
}

