package shared_init

import (
	"github.com/spacetimi/server/timi_shared/code/config"
	"github.com/spacetimi/server/timi_shared/code/core/adaptors/mongo_wrapper"
	"github.com/spacetimi/server/timi_shared/code/core/adaptors/redis_adaptor"
	"github.com/spacetimi/server/timi_shared/code/core/services/metadata_service"
)

func SharedInit(appInitializer IAppInitializer) {

	mongo_wrapper.Initialize(config.GetEnvironmentConfiguration().SharedMongoURL,
							 config.GetEnvironmentConfiguration().SharedDatabaseName,
							 config.GetEnvironmentConfiguration().AppMongoURL,
							 config.GetEnvironmentConfiguration().AppDatabaseName)

	redis_adaptor.Initialize()

	metadata_service.Initialize()

	// End shared init. Can now do App init

	appInitializer.AppInit()
}

