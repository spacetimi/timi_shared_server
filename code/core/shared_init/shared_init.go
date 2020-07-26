package shared_init

import (
	"strings"

	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/mongo_adaptor"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_factory"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func SharedInit(appInitializer IAppInitializer) {

	config.Initialize(appInitializer.AppName())

	var reqdServices *RequiredServices
	reqdServicesFilePath := config.GetAppConfigFilesPath() + "/services/required_services.json"
	err := config.ReadConfigFile(reqdServicesFilePath, &reqdServices)
	if err != nil {
		logger.LogFatal("error reading required-services config" +
			"|file path=" + reqdServicesFilePath +
			"|error=" + err.Error())
		return
	}
	for _, serviceName := range reqdServices.Services {
		initializeService(serviceName)
	}

	// End shared init. Can now do App init

	err = appInitializer.AppInit()
	if err != nil {
		logger.LogFatal("error during app init|error=" + err.Error())
	}
}

func initializeService(serviceName string) {
	logger.LogInfo("initializing service|service name=" + serviceName)

	switch serviceName {

	case "mongo_adaptor":
		configObject := mongo_adaptor.Config{}
		readConfigForService(serviceName, &configObject)
		mongo_adaptor.Initialize(configObject)

	case "redis_adaptor":
		configObject := redis_adaptor.Config{}
		readConfigForService(serviceName, &configObject)
		redis_adaptor.Initialize(configObject)

	case "metadata_service":
		metadata_service.Initialize()
		metadata_factory.Initialize()
		registerMetadataFactories()

	case "identity_service":
		identity_service.Initialize()

	default:
		logger.LogError("attempting to initialize unknown service" +
			"|service name=" + serviceName)
		return
	}
}

func readConfigForService(serviceName string, configObject interface{}) {
	configFilePath := config.GetAppConfigFilesPath() + "/services/" + serviceName + "/" +
		strings.ToLower(config.GetEnvironmentConfiguration().AppEnvironment.String()) + ".json"
	err := config.ReadConfigFile(configFilePath, &configObject)
	if err != nil {
		logger.LogFatal("error reading config for service" +
			"|service name=" + serviceName +
			"|file path=" + configFilePath +
			"|error=" + err.Error())
	}
}

// TODO: Avi: Move this somewhere else?
func registerMetadataFactories() {
	// Nothing yet
}
