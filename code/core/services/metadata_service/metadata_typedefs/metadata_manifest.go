package metadata_typedefs

type MetadataManifestItem struct {
	MetadataKey string
	Hash string
}

type MetadataManifest struct {
	MetadataManifestItems []MetadataManifestItem

	_itemsAsMap map[string]MetadataManifestItem 	// key => metadata manifest item for key
}

func (mm *MetadataManifest) Initialize() {
	mm._itemsAsMap = map[string]MetadataManifestItem{}
	for _, item := range mm.MetadataManifestItems {
		mm._itemsAsMap[item.MetadataKey] = item
	}
}

func (mm *MetadataManifest) GetManifestItem(key string) *MetadataManifestItem {
	if mm._itemsAsMap == nil {
		return nil
	}

	manifestItem, ok := mm._itemsAsMap[key]
	if !ok {
		return nil
	}

	return &manifestItem
}

