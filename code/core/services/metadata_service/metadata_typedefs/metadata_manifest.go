package metadata_typedefs

type MetadataManifestItem struct {
	MetadataKey string
	Hash string
}

type MetadataManifest struct {
	MetadataManifestItems []*MetadataManifestItem

	_itemsAsMap map[string]*MetadataManifestItem 	// key => metadata manifest item for key
}

func (mm *MetadataManifest) Initialize() {
	mm._itemsAsMap = map[string]*MetadataManifestItem{}
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

	return manifestItem
}

/**
 * Only meant to be called from the admin tool / scripts
 */
func (mm *MetadataManifest) SetManifestItem(key string, hash string) {
    for _, item := range mm.MetadataManifestItems {
    	if item.MetadataKey == key {
			item.Hash = hash
			return
		}
	}

    // Must be a new item
    manifestItem := &MetadataManifestItem{
    	MetadataKey:key,
    	Hash:hash,
	}
	mm.MetadataManifestItems = append(mm.MetadataManifestItems, manifestItem)
	mm._itemsAsMap[key] = manifestItem
}
