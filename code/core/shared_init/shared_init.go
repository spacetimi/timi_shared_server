package shared_init

import (
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_adaptor"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_factory"
)

func SharedInit(appInitializer IAppInitializer) {

	mongo_adaptor.Initialize(config.GetEnvironmentConfiguration().SharedMongoURL,
							 config.GetEnvironmentConfiguration().SharedDatabaseName,
							 config.GetEnvironmentConfiguration().AppMongoURL,
						  	 config.GetEnvironmentConfiguration().AppDatabaseName)

	redis_adaptor.Initialize()

	metadata_service.Initialize()
	metadata_factory.Initialize()

	identity_service.Initialize()

	registerMetadataFactories()

	// End shared init. Can now do App init

	appInitializer.AppInit()
}

// TODO: Avi: Move this somewhere else?
func registerMetadataFactories() {
    // Nothing yet
}

