package metadata_factory

import (
    "errors"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
)

var kMetadataFactories map[string]IMetadataFactory

func Initialize() {
    kMetadataFactories = make(map[string]IMetadataFactory, 0)
}

func RegisterFactory(metadataItemKey string, factory IMetadataFactory) {
    _, ok := kMetadataFactories[metadataItemKey]
    if ok {
        logger.LogWarning("trying to register duplicate factory|metadata item key=" + metadataItemKey)
        return
    }
    kMetadataFactories[metadataItemKey] = factory
}

func InstantiateMetadataItem(metadataItemKey string) (metadata_typedefs.IMetadataItem, error) {
    factory, ok := kMetadataFactories[metadataItemKey]
    if !ok {
        return nil, errors.New("no metadata item factory registered")
    }

    return factory.Instantiate(), nil
}

type IMetadataFactory interface {
    Instantiate() metadata_typedefs.IMetadataItem
}